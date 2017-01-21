package stack

import "log"
import "io/ioutil"

import "crypto/x509"
import "encoding/pem"
import "encoding/json"

import tlvCodec "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/tables/lpm"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/stack/config"
import "github.com/chris-wood/spud/stack/cache"
import "github.com/chris-wood/spud/stack/pit"
import "github.com/chris-wood/spud/stack/component/connector"
import "github.com/chris-wood/spud/stack/component/codec"
import "github.com/chris-wood/spud/stack/component/crypto"
import "github.com/chris-wood/spud/stack/component/crypto/validator"
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
    head Component
    pendingMap map[string]MessageCallback
    pendingQueue chan messages.MessageWrapper
    prefixTable lpm.LPM
}

func (s *Stack) Enqueue(msg messages.MessageWrapper) {
    s.head.Enqueue(msg)
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
    s.head.Enqueue(msg)
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
        msg := s.head.Dequeue()

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
    }
}

func Create(config config.StackConfig) (*Stack, error) {
  var err error

  // Create the link
  var fc connector.ForwarderConnector
  switch config.Link {
  case "tcp":
      locator := config.FwdAddress
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

  // Create the shared data structures
  stackCache := cache.NewCache()
  stackPit := pit.NewPIT()

  // Create codec
  codecComponent := codec.NewCodecComponent(fc, stackCache, stackPit)

  // Create crypto component
  // XXX: defer validator and encryptor creation to the crypto component
  cryptoContext := context.NewCryptoContext()
  cryptoComponent := crypto.NewCryptoComponent(cryptoContext, codecComponent)

  // Create the right crypto processors
  if len(config.Keys) > 0 {
      for _, keyFileName := range(config.Keys) {
          keyData, err := ioutil.ReadFile(keyFileName)
          block, _ := pem.Decode(keyData)
          privateKey, parseError := x509.ParsePKCS1PrivateKey(block.Bytes)
          if parseError != nil {
              log.Printf("Failed to parse private key: %s", err)
          } else {
              rsaProcessor, _ := validator.NewRSAProcessorWithKey(privateKey)
              cryptoComponent.AddCryptoProcessor("/", rsaProcessor)
          }
      }
  } else {
      rsaProcessor, _ := validator.NewRSAProcessor(2048)
      cryptoComponent.AddCryptoProcessor("/", rsaProcessor)
  }

  // Assemble the stack
  stack := &Stack{
      components: []Component{cryptoComponent, codecComponent},
      head: cryptoComponent,
      pendingQueue: make(chan messages.MessageWrapper, config.PendingBufferSize), // random constant -- make this configurable
      pendingMap: make(map[string]MessageCallback),
      prefixTable: lpm.LPM{},
  }

  // Start each component
  for _, component := range(stack.components) {
      go component.ProcessEgressMessages()
      go component.ProcessIngressMessages()
  }
  go stack.processInputQueue()

  return stack, nil
}

func CreateRaw(configString string) (*Stack, error) {
    var configStruct config.StackConfig
    err := json.Unmarshal([]byte(configString), &configStruct)
    if err == nil {
        return Create(configStruct)
    }
    return nil, err
}

func CreateTest() (*Stack, error) {
    return CreateRaw(`{"link": "loopback"}`)
}
