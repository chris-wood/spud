package stack

// import "fmt"

import tlvCodec "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/stack/cache"
import "github.com/chris-wood/spud/stack/pit"
import "github.com/chris-wood/spud/stack/component/connector"
import "github.com/chris-wood/spud/stack/component/codec"
import "github.com/chris-wood/spud/stack/component/crypto"
import "github.com/chris-wood/spud/stack/component/crypto/processor"
import "github.com/chris-wood/spud/stack/component/crypto/context"

type MessageCallback func(msg messages.Message)

type Component interface {
    // XXX: rename to Push and Pop, respectively
    Enqueue(messages.Message)
    Dequeue() messages.Message
    ProcessEgressMessages()
    ProcessIngressMessages()
}

type Stack struct {
    components []Component
    pendingMap map[string]MessageCallback

    pendingQueue chan messages.Message
}

func (s *Stack) Enqueue(msg messages.Message) {
    s.components[0].Enqueue(msg)
}

func (s *Stack) Dequeue() messages.Message {
    return <- s.pendingQueue
}

func (s *Stack) Get(msg messages.Message, callback MessageCallback) {
    if (msg.GetPacketType() != tlvCodec.T_INTEREST) {
        return
    }

    // Register a callback for this message
    // XXX

    // Enqueue the message into the top of the stack
    s.components[0].Enqueue(msg)
}

func (s *Stack) Put(msg messages.Message) {
    if (msg.GetPacketType() == tlvCodec.T_INTEREST) {
        return
    }
    s.components[0].Enqueue(msg)
}

func (s *Stack) processInputQueue() {
    for ;; {
        // Dequeue messages as they arrive
        msg := s.components[0].Dequeue()

        // Check to see if there's a pending callback
        // XXX

        // Nope -- enqueue the message in the pending queue to free up
        // space in the first component's channel
        // This will block until it's ready to do something
        s.pendingQueue <- msg
    }
}

// Get(message, callback)
// Put(message) -- there is no callback for this!

/*
{
    "connector" : "tcp"
    "keystore": "<path to key store>"
}
*/
func Create(config string) *Stack {
    // 1. create connector
    fc, _ := connector.NewLoopbackForwarderConnector()

    // 1.5. create the shared data structures
    stackCache := cache.NewCache()
    stackPit := pit.NewPIT()

    // 2. create codec
    codecComponent := codec.NewCodecComponent(fc, stackCache, stackPit)

    // 3. create crypto component
    // XXX: the processor information would be pulled from the configuration file
    // XXX: check crypto processor errors here
    rsaProcessor, _ := processor.NewRSAProcessor(2048)
    cryptoContext := context.NewCryptoContext()
    cryptoComponent := crypto.NewCryptoComponent(cryptoContext, codecComponent)
    cryptoComponent.AddCryptoProcessor("/", rsaProcessor)

    // 4. assemble the stack
    stack := &Stack{
        components: []Component{cryptoComponent, codecComponent},
        pendingQueue: make(chan messages.Message, 10000), // random constant -- make this configurable
    }

    // 5. start each component
    for _, component := range(stack.components) {
        go component.ProcessEgressMessages()
        go component.ProcessIngressMessages()
    }
    go stack.processInputQueue()

    return stack
}
