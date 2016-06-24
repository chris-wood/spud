package stack

import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/stack/connector"
import "github.com/chris-wood/spud/stack/codec"

type Component interface {
    UpstreamQueue() chan *messages.Message
    DownstreamQueue() chan *messages.Message
}

type Stack struct {
    components []Component

    stackCodec codec.Codec
    forwarderConnector connector.ForwarderConnector
}

/*
{
    "connector" : "tcp"
    "keystore": "<path to key store>"
}
*/
func Create(config string) *Stack {
    // 1. create connector
    fc, _ := connector.NewLoopbackForwarderConnector()

    // 2. create codec
    stackCodec := codec.NewCodec(fc)

    // 3. create other components
    // authenticator := crypto.XXXX

    // 4. assemble the stack
    return &Stack{components: nil, stackCodec: stackCodec, forwarderConnector: fc}
}
