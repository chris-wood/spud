package component

import "github.com/chris-wood/spud/messages"

type Component interface {
	Push(*messages.MessageWrapper)
	Inject(*messages.MessageWrapper)
	Pop() *messages.MessageWrapper
	ProcessEgressMessages()
	ProcessIngressMessages()
}
