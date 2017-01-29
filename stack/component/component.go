package component

import "github.com/chris-wood/spud/messages"

type Component interface {
	// XXX: rename to Push and Pop, respectively
	Enqueue(*messages.MessageWrapper)
	Dequeue() *messages.MessageWrapper
	ProcessEgressMessages()
	ProcessIngressMessages()
}
