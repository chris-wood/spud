package adapter

import "fmt"
import "time"

import "github.com/chris-wood/spud/stack/api/portal"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/messages/interest"
import "github.com/chris-wood/spud/messages/content"

type KVSAPI struct {
    p *portal.Portal
}

type KVSAPIError struct {
    arg  int
    prob string
}
func (e KVSAPIError) Error() string {
    return fmt.Sprintf("%d - %s", e.arg, e.prob)
}

type RequestCallback func(string, []byte) []byte
type ResponseCallback func([]byte)

func NewKVSAPI(thePortal *portal.Portal) *KVSAPI {
    api := &KVSAPI{
        p: thePortal,
    }
    return api
}

func (n *KVSAPI) Get(nameString string, timeout time.Duration) ([]byte, error) {
    requestName, err := name.Parse(nameString)
    if err == nil {
        request := messages.Package(interest.CreateWithName(requestName))
        response, err := n.p.Get(request, timeout)
        if err == nil {
            return response.Payload().Value(), nil
        } else {
            return []byte{}, err
        }
    }
    return nil, err
}

func (n *KVSAPI) GetAsync(nameString string, callback ResponseCallback) error {
    requestName, err := name.Parse(nameString)
    if err == nil {
        request := messages.Package(interest.CreateWithName(requestName))
        n.p.GetAsync(request, func(msg messages.MessageWrapper) {
            callback(msg.Payload().Value())
        })
    }
    return err
}

func (n *KVSAPI) Serve(nameString string, callback RequestCallback) {
    n.p.Serve(nameString, func(msg messages.MessageWrapper) messages.MessageWrapper {
        encapPayload := msg.Payload().Value()
        data := callback(msg.Identifier(), encapPayload)
        dataPayload := payload.Create(data)
        response := messages.Package(content.CreateWithNameAndPayload(msg.Name(), dataPayload))
        return response
    })
}
