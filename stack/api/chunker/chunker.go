package chunker

import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/name"
// import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/messages/interest"
// import "github.com/chris-wood/spud/messages/content"

type ChunkerAPI struct {
    apiStack *stack.Stack
}

type RequestCallback func(string, []byte) []byte
type ResponseCallback func([]byte)

func NewChunkerAPI(apiStack *stack.Stack) *ChunkerAPI {
    api := &ChunkerAPI{
        apiStack: apiStack,
    }

    return api
}

func (n *ChunkerAPI) Get(nameString string, callback ResponseCallback) {
    requestName, err := name.Parse(nameString)
    if err == nil {
        request := interest.CreateWithName(requestName)
        n.apiStack.Get(request, func(msg messages.Message) {
            // XXX: implement the chunking protocol
        })
    }
}

func (n *ChunkerAPI) Serve(nameString string, callback RequestCallback) {
    n.apiStack.Service(nameString, func(msg messages.Message) {
        // encapPayload := msg.Payload().Value()
        // data := callback(msg.Identifier(), encapPayload)
        // dataPayload := payload.Create(data)
        // response := content.CreateWithNameAndPayload(msg.Name(), dataPayload)
        // n.apiStack.Enqueue(response)
    })
}
