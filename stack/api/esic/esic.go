package esic

import "github.com/chris-wood/spud/tables/lpm"
import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/messages/interest"
import "github.com/chris-wood/spud/messages/content"

type ESIC struct {
    writeEncKey []byte
    writeMacKey []byte
    readEncKey []byte
    readMacKey []byte

    prefixTable lpm.LPM
    pendingMap map[string]ResponseCallback
    apiStack stack.Stack
}

type RequestCallback func(string, []byte) []byte
type ResponseCallback func([]byte)

func (n *ESIC) Serve(prefix string, sessionID string, callback RequestCallback) {
    // XX store the prefix, callback tuple in the map
    // n.prefixMap[prefix] = callback
    prefixName, err := name.Parse(prefix)
    if err == nil {
        nameComponents := prefixName.SegmentStrings()
        n.prefixTable.Insert(nameComponents, callback)
    }
}

func (n *ESIC) Get(nameString string, callback ResponseCallback) {
    requestName, err := name.Parse(nameString)
    if err == nil {
        request := interest.CreateWithName(requestName)

        // XXX: encode the interest
        // XXX: encrypt it
        // XXX: insert in interest payload
        // XXX: send

        n.pendingMap[request.Identifier()] = callback
        n.apiStack.Enqueue(request)
    }
}

func (n *ESIC) process() {
    for ;; {
        msg := n.apiStack.Dequeue()

        if msg.GetPacketType() == codec.T_INTEREST {
            // XXX: this needs to use LPM to identify the service prefix
            requestName := msg.Name()
            numSsegments := len(requestName.Segments)

            for index := 1; index <= numSsegments; index++ {
                prefix := requestName.Prefix(index)
                nameComponents := requestName.SegmentStrings()
                // callback, ok := n.prefixMap[prefix]
                callbackInterface, ok := n.prefixTable.Lookup(nameComponents)

                if ok {
                    callback := callbackInterface.(RequestCallback)
                    go func() {
                        // XXX: this needs to be the full name string, not the prefix
                        data := callback(prefix, msg.Payload().Value())
                        dataPayload := payload.Create(data)
                        response := content.CreateWithNameAndPayload(requestName, dataPayload)
                        n.apiStack.Enqueue(response)
                    }()
                    break
                }
            }
        } else {
            callback, ok := n.pendingMap[msg.Identifier()]
            if ok {
                pay := msg.Payload()
                callback(pay.Value())
            }
        }

        // extract the name from the message hand it to the callback
        // enqueue the message to the stack
    }
}
