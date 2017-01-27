package esic

// import "github.com/chris-wood/spud/stack"
// import "github.com/chris-wood/spud/codec"
// import "github.com/chris-wood/spud/messages"
// import "github.com/chris-wood/spud/messages/name"
// import "github.com/chris-wood/spud/messages/payload"
// import "github.com/chris-wood/spud/messages/interest"
// import "github.com/chris-wood/spud/messages/content"

import "crypto/aes"
import "crypto/cipher"

type ESIC struct {
    sessionID string
    counter int

    writeEncKey []byte
    writeMacKey []byte
    readEncKey []byte
    readMacKey []byte

    WriteCipher cipher.AEAD
    ReadCipher cipher.AEAD
}

// type RequestCallback func(string, []byte) []byte
// type ResponseCallback func([]byte)

// XXX: session creation

func NewESIC(masterSecret []byte, sessionID string) *ESIC {
    WriteCipher, err := aes.NewCipher(masterSecret)
    if err != nil {
        panic(err.Error())
    }
    writeAEAD, err := cipher.NewGCM(WriteCipher)
    if err != nil {
        panic(err.Error())
    }

    ReadCipher, err := aes.NewCipher(masterSecret)
    if err != nil {
        panic(err.Error())
    }
    readAEAD, err := cipher.NewGCM(ReadCipher)
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
        WriteCipher: writeAEAD,
        ReadCipher: readAEAD,
    }

    return &esic
}

// XXX: add functions for encryption and decryption

// func (n *ESIC) Serve(prefix string, callback RequestCallback) {
//     n.apiStack.Service(prefix, func(msg messages.MessageWrapper) {
//         encryptedPayload := msg.InnerMessage().Payload().Value()
//
//         // Decrypt the interest
//         nonce := make([]byte, 12) // all zeros to start
//         encodedInterest, err := n.ReadCipher.Open(nil, nonce, encryptedPayload, nil)
//         if err != nil {
//             return
//         }
//
//         d := codec.Decoder{}
//         decodedTlV := d.Decode(encodedInterest)
//         if len(decodedTlV) > 1 {
//             return // there should be only one thing encapsulated -- a piece of data
//         }
//         responseMsg, err := messages.CreateFromTLV(decodedTlV)
//         if err != nil {
//             return  // handle this as needed
//         }
//
//         data := callback(prefix, responseMsg.Payload().Value())
//         dataPayload := payload.Create(data)
//         response := messages.ContentWrapper(content.CreateWithNameAndPayload(msg.Name(), dataPayload))
//
//         // Encode and encrypt the content response
//         // e := codec.Encoder{}
//         // encodedResponse := e.EncodeTLV(response)
//         encodedResponse := response.Encode()
//
//         // Encrypt the response
//         nonce = make([]byte, 12) // all zeros to start
//         encryptedResponse := n.WriteCipher.Seal(nil, nonce, encodedResponse, nil)
//
//         // Encap the response and send it downstream
//         wrappedPayload := payload.Create(encryptedResponse)
//         encapResponse := messages.ContentWrapper(content.CreateWithNameAndTypedPayload(msg.Name(), codec.T_PAYLOADTYPE_ENCAP, wrappedPayload))
//
//         n.apiStack.Enqueue(encapResponse)
//     })
// }

// func (n *ESIC) Get(nameString string, callback ResponseCallback) {
//     requestName, err := name.Parse(nameString)
//     if err == nil {
//         request := messages.InterestWrapper(interest.CreateWithName(requestName))
//
//         // e := codec.Encoder{}
//         // encodedRequest := e.EncodeTLV(request)
//         encodedRequest := request.Encode()
//
//         // Encrypt the interest
//         nonce := make([]byte, 12) // all zeros to start
//         encryptedInterest := n.WriteCipher.Seal(nil, nonce, encodedRequest, nil)
//
//         baseName, _ := name.Parse(nameString)
//         sessionName, _ := baseName.AppendComponent(n.sessionID)
//
//         encapPayload := payload.Create(encryptedInterest)
//         encapInterest := messages.InterestWrapper(interest.CreateWithNameAndPayload(sessionName, codec.T_PAYLOADTYPE_ENCAP, encapPayload))
//
//         n.apiStack.Get(encapInterest, func(msg messages.MessageWrapper) {
//             encryptedPayload := msg.InnerMessage().Payload().Value()
//
//             // Decrypt the interest
//             nonce := make([]byte, 12) // all zeros to start
//             encapInterest, err := n.ReadCipher.Open(nil, nonce, encryptedPayload, nil)
//             if err != nil {
//                 return
//             }
//
//             d := codec.Decoder{}
//             decodedTlV := d.Decode(encapInterest)
//             if len(decodedTlV) > 1 {
//                 return // there should be only one thing encapsulated -- a piece of data
//             }
//
//             responseMsg, err := messages.CreateFromTLV(decodedTlV)
//             if err != nil {
//                 return  // handle this as needed
//             }
//
//             callback(responseMsg.InnerMessage().Payload().Value())
//         })
//     }
// }
