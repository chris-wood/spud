package adapter

import "fmt"
import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/interest"

type NameAPI struct {
    apiStack stack.Stack
    prefixMap map[string]RequestCallback
}

type RequestCallback func(string) []byte
type ResponseCallback func([]byte)

func NewNameAPI(s stack.Stack) NameAPI {
    api := NameAPI{apiStack: s, prefixMap: make(map[string]RequestCallback)}
    go api.process()
    return api
}

func (n NameAPI) Get(nameString string, callback ResponseCallback) {
    requestName, err := name.Parse(nameString)
    if err == nil {
        request := interest.CreateWithName(requestName)
        n.apiStack.Enqueue(request)
    }
}

func (n NameAPI) process() {
    for ;; {
        msg := n.apiStack.Dequeue()
        fmt.Println(msg.Identifier())
        // extract the name from the message hand it to the callback
        // enqueue the message to the stack
    }
}

func (n NameAPI) Serve(prefix string, callback RequestCallback) {
    // XX store the prefix, callback tuple in the map
}
