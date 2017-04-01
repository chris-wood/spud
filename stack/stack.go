package stack

import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/stack/component/tunnel"

type MessageCallback func(msg *messages.MessageWrapper)

type Stack interface {
    Enqueue(msg *messages.MessageWrapper)
    Dequeue() *messages.MessageWrapper
    Cancel(msg *messages.MessageWrapper)
    Get(msg *messages.MessageWrapper, callback MessageCallback)
    Service(prefix *name.Name, callback MessageCallback)
    AddSession(session *tunnel.Session, baseName *name.Name)
}
