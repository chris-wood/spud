package esic

import "github.com/chris-wood/spud/tables/lpm"
import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/messages/interest"
import "github.com/chris-wood/spud/messages/content"

type ESIC struct {
    sessionID string
    counter int

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

// XXX: session creation

func NewESIC(masterSecret []byte, sessionID string) *ESIC {
    esic := ESIC{
        sessionID: sessionID,
        counter: 0,
        writeMacKey: masterSecret,
        writeEncKey: masterSecret,
        readMacKey: masterSecret,
        readEncKey: masterSecret,
    }

    // XXX: invoke the process function here

    return &esic
}

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

        e := codec.Encoder{}
        encodedRequest := e.EncodeTLV(request)
        // XXX: encrypt+MAC the interest here

        baseName, _ := name.Parse(nameString)
        sessionName, _ := baseName.AppendComponent(n.sessionID)

        encapPayload := payload.Create(encodedRequest)
        encapInterest := interest.CreateWithNameAndPayload(sessionName, codec.T_PAYLOADTYPE_ENCAP, encapPayload)

        n.pendingMap[encapInterest.Identifier()] = callback
        // n.apiStack.Enqueue(encapInterest)

        n.apiStack.Get(encapInterest, func(msg messages.Message) {
            encapPayload := msg.Payload().Value()
            // XXX: decrypt the payload

            d := codec.Decoder{}
            decodedTlV := d.Decode(encapPayload)
            if len(decodedTlV) > 1 {
                return // there should be only one thing encapsulated -- a piece of data
            }
            responseMsg, err := messages.CreateFromTLV(decodedTlV)
            if err != nil {
                return  // handle this as needed
            }

            callback(responseMsg.Payload().Value())
        })
    }
}

// func (n *ESIC) processResponse(msg messages.Message) {
//
// }

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
                        encapPayload := msg.Payload().Value()
                        // XXX: decrypt the payload

                        d := codec.Decoder{}
                        decodedTlV := d.Decode(encapPayload)
                        if len(decodedTlV) > 1 {
                            return // there should be only one thing encapsulated -- a piece of data
                        }
                        responseMsg, err := messages.CreateFromTLV(decodedTlV)
                        if err != nil {
                            return  // handle this as needed
                        }

                        data := callback(prefix, responseMsg.Payload().Value())
                        dataPayload := payload.Create(data)
                        response := content.CreateWithNameAndPayload(requestName, dataPayload)

                        // Encode and encrypt the content response
                        e := codec.Encoder{}
                        encodedResponse := e.EncodeTLV(response)
                        // XXX: encrypt the response

                        // Encap the response and send it downstream
                        wrappedPayload := payload.Create(encodedResponse)
                        encapResponse := content.CreateWithNameAndTypedPayload(msg.Name(), codec.T_PAYLOADTYPE_ENCAP, wrappedPayload)

                        n.apiStack.Enqueue(encapResponse)
                    }()
                    break
                }
            }
        }

        // extract the name from the message hand it to the callback
        // enqueue the message to the stack
    }
}
