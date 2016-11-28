package stack

import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/stack/cache"
import "github.com/chris-wood/spud/stack/pit"
import "github.com/chris-wood/spud/stack/component/connector"
import "github.com/chris-wood/spud/stack/component/codec"
import "github.com/chris-wood/spud/stack/component/crypto"
import "github.com/chris-wood/spud/stack/component/crypto/processor"
import "github.com/chris-wood/spud/stack/component/crypto/context"

type Component interface {
    Enqueue(messages.Message)
    Dequeue() messages.Message
    ProcessEgressMessages()
    ProcessIngressMessages()
}

type Stack struct {
    components []Component
}

func (s Stack) Enqueue(msg messages.Message) {
    s.components[0].Enqueue(msg)
}

func (s Stack) Dequeue() messages.Message {
    return s.components[0].Dequeue()
}

/*
{
    "connector" : "tcp"
    "keystore": "<path to key store>"
}
*/
func Create(config string) Stack {
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
    stack := Stack{components: []Component{cryptoComponent, codecComponent}}

    // 5. start each component
    for _, component := range(stack.components) {
        go component.ProcessEgressMessages()
        go component.ProcessIngressMessages()
    }

    return stack
}
