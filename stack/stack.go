package stack

import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/stack/connector"
import "github.com/chris-wood/spud/stack/codec"
import "github.com/chris-wood/spud/stack/crypto"
import "github.com/chris-wood/spud/stack/crypto/processor"

type Component interface {
    Enqueue(messages.Message)
    Dequeue() messages.Message
    ProcessEgressMessages()
    ProcessIngressMessages()
}

type Stack struct {
    components []Component
    codecComponent codec.CodecComponent
    forwarderConnector connector.ForwarderConnector
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

    // 2. create codec
    codecComponent := codec.NewCodecComponent(fc)
    go codecComponent.ProcessEgressMessages()
    go codecComponent.ProcessIngressMessages()

    // 3. create crypto component
    // XXX: the processor information would be pulled from the configuration file
    // XXX: check crypto processor errors here
    rsaProcessor, _ := processor.NewRSAProcessor(2048)
    cryptoComponent := crypto.NewCryptoComponent(rsaProcessor, codecComponent)
    go cryptoComponent.ProcessEgressMessages()
    go cryptoComponent.ProcessIngressMessages()

    // 4. assemble the stack
    return Stack{components: []Component{cryptoComponent, codecComponent}, codecComponent: codecComponent, forwarderConnector: fc}
}
