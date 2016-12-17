package messages

import "fmt"
import "hash"
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
    // Messages can encode themselves
    Encode() []byte

    // Messages have names, identifiers, and optionally, a payload
    Name() name.Name
    Identifier() string
    NamelessIdentifier() string
    Payload() payload.Payload
    PayloadType() uint8

    // Messages can compute the hashes of their protected regions and their complete packet formats.
    HashProtectedRegion(hasher hash.Hash) []byte
    ComputeMessageHash(hasher hash.Hash) []byte

    // Messages have validation information that are set outside of the messages themselves
    SetValidationAlgorithm(va validation.ValidationAlgorithm)
    SetValidationPayload(va validation.ValidationPayload)
    GetValidationAlgorithm() validation.ValidationAlgorithm
    GetValidationPayload() validation.ValidationPayload

    // Generic slot for containers
    AddContainer(container codec.TLV)
    GetContainer(containerType uint16) (codec.TLV, error)

    // Type APIs
    GetPacketType() uint16
    // GetMessageType() int
}

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

    return result, err
}
