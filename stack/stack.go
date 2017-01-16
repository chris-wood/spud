package stack

// import "fmt"
import "log"
import "ioutil"

import "encoding/pem"
import "encoding/json"

import tlvCodec "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/tables/lpm"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/stack/cache"
import "github.com/chris-wood/spud/stack/pit"
import "github.com/chris-wood/spud/stack/component/connector"
import "github.com/chris-wood/spud/stack/component/codec"
import "github.com/chris-wood/spud/stack/component/crypto"
import "github.com/chris-wood/spud/stack/component/crypto/processor"
import "github.com/chris-wood/spud/stack/component/crypto/context"

type MessageCallback func(msg messages.MessageWrapper)

type Component interface {
    // XXX: rename to Push and Pop, respectively
    Enqueue(messages.MessageWrapper)
    Dequeue() messages.MessageWrapper
    ProcessEgressMessages()
    ProcessIngressMessages()
}

type Stack struct {
    components []Component
    pendingMap map[string]MessageCallback
    pendingQueue chan messages.MessageWrapper
    prefixTable lpm.LPM
}

func (s *Stack) Enqueue(msg messages.MessageWrapper) {
    s.components[0].Enqueue(msg)
}

func (s *Stack) Dequeue() messages.MessageWrapper {
    return <- s.pendingQueue
}

func (s *Stack) Get(msg messages.MessageWrapper, callback MessageCallback) {
    if (msg.GetPacketType() != tlvCodec.T_INTEREST) {
        return
    }

    // Register the callback
    s.pendingMap[msg.Identifier()] = callback

    // Enqueue the message into the top of the stack
    s.components[0].Enqueue(msg)
}

func (s *Stack) Service(nameString string, callback MessageCallback) {
    prefixName, err := name.Parse(nameString)
    if err == nil {
        nameComponents := prefixName.SegmentStrings()
        s.prefixTable.Insert(nameComponents, callback)
    }
}

func (s *Stack) processInputQueue() {
    for ;; {
        // Dequeue messages as they arrive
        msg := s.components[0].Dequeue()

        switch (msg.GetPacketType()) {
        case tlvCodec.T_INTEREST:
            requestName := msg.Name()
            numSsegments := len(requestName.Segments)

            handled := false
            for index := 1; index <= numSsegments; index++ {
                nameComponents := requestName.SegmentStrings()
                callbackInterface, ok := s.prefixTable.Lookup(nameComponents)

                if ok {
                    callback := callbackInterface.(MessageCallback)
                    callback(msg)
                    s.prefixTable.Drop(nameComponents)
                    handled = true
                }
            }

            if !handled {
                s.pendingQueue <- msg
            }
        default:
            callback, ok := s.pendingMap[msg.Identifier()]
            if ok {
                callback(msg)
            } else {
                s.pendingQueue <- msg
            }
        }

        // // Check to see if there's a pending callback
        // callback, ok := s.pendingMap[msg.Identifier()]
        // if ok && msg.GetPacketType() != tlvCodec.T_INTEREST {
        //     fmt.Println("Calling the response callback")
        //     callback(msg)
        // } else {
        //     // Nope -- enqueue the message in the pending queue to free up
        //     // space in the first component's channel
        //     // This will block until it's ready to do something
        //     s.pendingQueue <- msg
        // }
    }
}

/*
{
    "connector" : "tcp"
    "keystore": "<path to key store>"
}
*/
func Create(config string) *Stack {
    // Decode the JSON config
    var configMap map[string]interface{}
    if err := json.Unmarshal([]byte(config), &configMap); err != nil {
        panic(err)
    }

    var err error

    // 1. create the link
    var fc connector.ForwarderConnector
    switch configMap["link"].(string) {
    case "tcp":
        locator := configMap["fwd-address"].(string)
        fc, err = connector.NewAthenaTCPForwarderConnector(locator)
        if err != nil {
            panic(err)
        }
        break
    case "loopback":
        fallthrough
    default:
        fc, _ = connector.NewLoopbackForwarderConnector()
        break
    }

    if fc == nil {
        log.Panic("Catastrophic failure: a connector was not created.")
    }

    // 1.5. create the shared data structures
    stackCache := cache.NewCache()
    stackPit := pit.NewPIT()

    // 2. create codec
    codecComponent := codec.NewCodecComponent(fc, stackCache, stackPit)

    // 3. create crypto component
    cryptoContext := context.NewCryptoContext()
    cryptoComponent := crypto.NewCryptoComponent(cryptoContext, codecComponent)

    // 3.5. create the right crypto processors
    keyList, ok := configMap["keys"]
    if ok {
        for _, keyFileName := range(keyList.(string)) {
            keyData, err := ioutil.ReadFile(keyFileName)
            block := pem.Decode(keyData)
            privateKey, parseError := ParsePKCS1PrivateKey(block.Bytes)
            if parseError != nil {
                log.Printf("Failed to parse private key: %s", err)
            } else {
                rsaProcessor, _ := NewRSAProcessorWithKey(privateKey)
                cryptoComponent.AddCryptoProcessor("/", rsaProcessor)
            }
        }
    } else {
        rsaProcessor, _ := processor.NewRSAProcessor(2048)
        cryptoComponent.AddCryptoProcessor("/", rsaProcessor)
    }

    // 4. assemble the stack
    stack := &Stack{
        components: []Component{cryptoComponent, codecComponent},
        pendingQueue: make(chan messages.MessageWrapper, 10000), // random constant -- make this configurable
        pendingMap: make(map[string]MessageCallback),
        prefixTable: lpm.LPM{},
    }

    // 5. start each component
    for _, component := range(stack.components) {
        go component.ProcessEgressMessages()
        go component.ProcessIngressMessages()
    }
    go stack.processInputQueue()

    return stack
}

func CreateTest() *Stack {
    return Create(`{"link": "loopback"}`)
}
