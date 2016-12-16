package ccnxke

import "github.com/chris-wood/spud/tables/lpm"
import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/messages/name"
// import "github.com/chris-wood/spud/messages/payload"
// import "github.com/chris-wood/spud/messages/interest"
// import "github.com/chris-wood/spud/messages/content"

type CCNxKEAPI struct {
    apiStack stack.Stack

    prefixTable lpm.LPM
    pendingMap map[string]ResponseCallback
}

type RequestCallback func(string, []byte) []byte
type ResponseCallback func([]byte)

func NewCCNxKEAPI(s stack.Stack) *CCNxKEAPI {
    // prefixLPM := lpm.LPM{}
    // api := &NameAPI{
    //     apiStack: s,
    //     prefixTable: prefixLPM,
    //     pendingMap: make(map[string]ResponseCallback),
    // }
    //
    // go api.process()
    // return api

    return nil
}

func (n *CCNxKEAPI) Connect(prefix name.Name) {
    // XXX: do the handshake establishment here...

}

func (n *CCNxKEAPI) Service(prefix name.Name) bool {
    // XXX: advertise route for this prefix and manage session state
    return false
}

func (n *CCNxKEAPI) process() {
    for ;; {
        // msg := n.apiStack.Dequeue()
        // XXX: process the message, enqueue or dequeue as needed
    }
}
