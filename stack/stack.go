package stack

import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/stack/connector"
import "github.com/chris-wood/spud/stack/codec"

type Component interface {
    Enqueue(messages.Message)
    Dequeue() messages.Message
    ProcessEgressMessages()
    ProcessIngressMessages()
}

type Stack struct {
    components []Component
    stackCodec codec.CodecComponent
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
    stackCodec := codec.NewCodecComponent(fc)
    go stackCodec.ProcessEgressMessages()
    go stackCodec.ProcessIngressMessages()

    // 3. create other components
    // authenticator := crypto.XXXX

    // 4. assemble the stack
    return Stack{components: []Component{stackCodec}, stackCodec: stackCodec, forwarderConnector: fc}
}
