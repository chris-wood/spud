package portal

import "fmt"
import "time"

import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/messages"

type Portal struct {
    apiStack *stack.Stack
}

type PortalError struct {
    arg  int
    prob string
}

func (e PortalError) Error() string {
    return fmt.Sprintf("%d - %s", e.arg, e.prob)
}

type RequestMessageCallback func(messages.MessageWrapper) messages.MessageWrapper
type ResponseMessageCallback func(messages.MessageWrapper)

func NewPortal(s *stack.Stack) *Portal {
    api := &Portal{
        apiStack: s,
    }

    return api
}

func (n *Portal) Get(request messages.MessageWrapper, timeout time.Duration) (messages.MessageWrapper, error) {
    signalChannel := make(chan messages.MessageWrapper, 1)
    n.apiStack.Get(request, func(msg messages.MessageWrapper) {
        signalChannel <- msg
    })

    var response messages.MessageWrapper
    select {
    case data := <- signalChannel:
        return data, nil
    case <-time.After(timeout):
        return response, PortalError{0, "Timeout"}
    }
}

func (n *Portal) GetAsync(request messages.MessageWrapper, callback ResponseMessageCallback) {
    n.apiStack.Get(request, func(msg messages.MessageWrapper) {
        callback(msg)
    })
}

func (n *Portal) Serve(nameString string, callback RequestMessageCallback) {
    n.apiStack.Service(nameString, func(msg messages.MessageWrapper) {
        response := callback(msg)
        n.apiStack.Enqueue(response)
    })
}
