package messages

import "fmt"
import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages/validation"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/messages/content"
import "github.com/chris-wood/spud/messages/interest"

type messageError struct {
    prob string
}

func (e messageError) Error() string {
    return fmt.Sprintf("%s", e.prob)
}

type Message interface {
    Name() name.Name
    Identifier() string
    HashSensitiveRegion() []byte
    ComputeMessageHash() []byte
    Encode() []byte
    Payload() payload.Payload

    SetValidationAlgorithm(va validation.ValidationAlgorithm)

    // XXX: should this be `get message type`?
    IsRequest() bool

    // XXX: this should take a signer as input
    // TagAndEncode() []byte
}

// XXX: create the right type of message here
func CreateFromTLV(tlv []codec.TLV) (Message, error) {
    var result Message
    var err error

    root := tlv[0]
    switch (root.Type()) {
    case codec.T_INTEREST:
        result, err = interest.CreateFromTLV(root)
        break
    case codec.T_OBJECT:
        result, err = content.CreateFromTLV(root)
        break
    default:
        fmt.Println("invalid type " + string(root.Type()))
        fmt.Println(root.Type())
        return result, messageError{"Unable to create a message from the top-level TLV type " + string(root.Type())}
    }

    if err != nil {
        fmt.Println("tried and failed to create a message from a TLV")
    } else {
        fmt.Println("Reconstructed message: " + result.Identifier())
    }

    return result, err
}
