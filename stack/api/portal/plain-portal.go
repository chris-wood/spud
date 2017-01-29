package portal

import "time"

import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/stack"

type PlainPortal struct {
	apiStack stack.Stack
}

func NewPortal(s stack.Stack) PlainPortal {
	api := PlainPortal{
		apiStack: s,
	}

	return api
}

func (n PlainPortal) Get(request *messages.MessageWrapper, timeout time.Duration) (*messages.MessageWrapper, error) {
	signalChannel := make(chan *messages.MessageWrapper, 1)
	n.apiStack.Get(request, func(msg *messages.MessageWrapper) {
		signalChannel <- msg
	})

	select {
	case data := <-signalChannel:
		return data, nil
	case <-time.After(timeout):
		return nil, PortalError{0, "Timeout"}
	}
}

func (n PlainPortal) GetAsync(request *messages.MessageWrapper, callback ResponseMessageCallback) {
	n.apiStack.Get(request, func(msg *messages.MessageWrapper) {
		callback(msg)
	})
}

func (n PlainPortal) Serve(prefix name.Name, callback RequestMessageCallback) {
	n.apiStack.Service(prefix, func(msg *messages.MessageWrapper) {
		response := callback(msg)
		n.apiStack.Enqueue(response)
	})
}
