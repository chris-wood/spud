package adapter

import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/messages/interest"
import "github.com/chris-wood/spud/messages/content"

type NameAPI struct {
    apiStack *stack.Stack
}

type RequestCallback func(string, []byte) []byte
type ResponseCallback func([]byte)

func NewNameAPI(s *stack.Stack) *NameAPI {
    api := &NameAPI{
        apiStack: s,
    }

    return api
}

func (n *NameAPI) Get(nameString string, callback ResponseCallback) {
    requestName, err := name.Parse(nameString)
    if err == nil {
        request := messages.InterestWrapper(interest.CreateWithName(requestName))
        n.apiStack.Get(request, func(msg messages.MessageWrapper) {
            callback(msg.InnerMessage().Payload().Value())
        })
    }
}

func (n *NameAPI) Serve(nameString string, callback RequestCallback) {
    n.apiStack.Service(nameString, func(msg messages.MessageWrapper) {
        encapPayload := msg.Payload().Value()
        data := callback(msg.Identifier(), encapPayload)
        dataPayload := payload.Create(data)
        response := messages.ContentWrapper(content.CreateWithNameAndPayload(msg.Name(), dataPayload))
        n.apiStack.Enqueue(response)
    })
}
