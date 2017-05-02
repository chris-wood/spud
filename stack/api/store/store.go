package store

import "fmt"
import "time"

import "github.com/chris-wood/spud/stack/api/portal"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/messages/interest"
import "github.com/chris-wood/spud/messages/content"

type StoreAPI struct {
	p portal.Portal
}

type StoreAPIError struct {
	arg  int
	prob string
}

func (e StoreAPIError) Error() string {
	return fmt.Sprintf("%d - %s", e.arg, e.prob)
}

type RequestCallback func(string, []byte) []byte
type ResponseCallback func([]byte)

func NewStoreAPI(thePortal portal.Portal) *StoreAPI {
	api := &StoreAPI{
		p: thePortal,
	}
	return api
}

func (n *StoreAPI) Get(nameString string, timeout time.Duration) ([]byte, error) {
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

func (n *StoreAPI) GetAsync(nameString string, callback ResponseCallback) error {
	requestName, err := name.Parse(nameString)
	if err == nil {
		request := messages.Package(interest.CreateWithName(requestName))
		n.p.GetAsync(request, func(msg *messages.MessageWrapper) {
			callback(msg.Payload().Value())
		})
	}
	return err
}

func (n *StoreAPI) Serve(nameString string, callback RequestCallback) {
	prefix, err := name.Parse(nameString)
	if err == nil {
		n.p.Serve(prefix, func(msg *messages.MessageWrapper) *messages.MessageWrapper {
			encapPayload := msg.Payload().Value()
			data := callback(msg.Identifier(), encapPayload)
			dataPayload := payload.Create(data)
			response := messages.Package(content.CreateWithNameAndPayload(msg.Name(), dataPayload))
			return response
		})
	}
}
