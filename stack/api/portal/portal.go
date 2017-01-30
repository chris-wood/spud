package portal

import "fmt"
import "time"

import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/name"

type RequestMessageCallback func(*messages.MessageWrapper) *messages.MessageWrapper
type ResponseMessageCallback func(*messages.MessageWrapper)

type Portal interface {
	Get(request *messages.MessageWrapper, timeout time.Duration) (*messages.MessageWrapper, error)
    GetAsync(request *messages.MessageWrapper, callback ResponseMessageCallback)
    Serve(prefix name.Name, callback RequestMessageCallback)
    Produce(data *messages.MessageWrapper)
}

type PortalError struct {
	Arg  int
	Prob string
}

func (e PortalError) Error() string {
	return fmt.Sprintf("%d - %s", e.Arg, e.Prob)
}
