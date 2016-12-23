package esic

import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/messages/interest"
import "github.com/chris-wood/spud/messages/content"

import "crypto/aes"
import "crypto/cipher"

type ESIC struct {
    sessionID string
    counter int

    writeEncKey []byte
    writeMacKey []byte
    readEncKey []byte
    readMacKey []byte

    writeCipher cipher.Block
    readCipher cipher.Block

    apiStack *stack.Stack
}

type RequestCallback func(string, []byte) []byte
type ResponseCallback func([]byte)

// XXX: session creation

func NewESIC(stack *stack.Stack, masterSecret []byte, sessionID string) *ESIC {
    writeCipher, err := aes.NewCipher(masterSecret)
    if err != nil {
        panic(err.Error())
    }

    readCipher, err := aes.NewCipher(masterSecret)
    if err != nil {
        panic(err.Error())
    }

    esic := ESIC{
        sessionID: sessionID,
        counter: 0,
        writeMacKey: masterSecret,
        writeEncKey: masterSecret,
        readMacKey: masterSecret,
        readEncKey: masterSecret,
        writeCipher: writeCipher,
        readCipher: readCipher,
        apiStack: stack,
    }

    return &esic
}

func (n *ESIC) Serve(prefix string, callback RequestCallback) {
    n.apiStack.Service(prefix, func(msg messages.Message) {
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
        response := content.CreateWithNameAndPayload(msg.Name(), dataPayload)

        // Encode and encrypt the content response
        e := codec.Encoder{}
        encodedResponse := e.EncodeTLV(response)
        // XXX: encrypt the response

        // Encap the response and send it downstream
        wrappedPayload := payload.Create(encodedResponse)
        encapResponse := content.CreateWithNameAndTypedPayload(msg.Name(), codec.T_PAYLOADTYPE_ENCAP, wrappedPayload)

        // n.apiStack.Enqueue(encapResponse)
        n.apiStack.Enqueue(encapResponse)
    })
}

func (n *ESIC) Get(nameString string, callback ResponseCallback) {
    requestName, err := name.Parse(nameString)
    if err == nil {
        request := interest.CreateWithName(requestName)

        e := codec.Encoder{}
        encodedRequest := e.EncodeTLV(request)

        /*
        aesgcm, err := cipher.NewGCM(n.writeCipher)
        if err != nil {
            panic(err.Error())
        }
        */
        // encryptedInterest := aesgcm.Seal(nil, nonce, encodedRequest, nil)

        // XXX: encrypt+MAC the interest here

        baseName, _ := name.Parse(nameString)
        sessionName, _ := baseName.AppendComponent(n.sessionID)

        encapPayload := payload.Create(encodedRequest)
        encapInterest := interest.CreateWithNameAndPayload(sessionName, codec.T_PAYLOADTYPE_ENCAP, encapPayload)

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

