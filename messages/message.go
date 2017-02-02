package messages

import "fmt"
import "hash"
import "crypto/sha256"
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

type MessageWrapper struct {
    msg Message
    validationAlgorithm validation.ValidationAlgorithm
    validationPayload validation.ValidationPayload
}

type Message interface {
    // Messages can encode themselves
    Encode() []byte

    // Messages have names, identifiers, and optionally, a payload
    Name() *name.Name
    Payload() *payload.Payload
    PayloadType() uint16

    // Identifier() string
    // NamelessIdentifier() string

    // Generic slots for containers
    AddContainer(container codec.TLV)
    GetContainer(containerType uint16) (codec.TLV, error)

    // Type APIs
    GetPacketType() uint16
    // GetMessageType() int
}

// We just need a single package function here
func Package(m Message) *MessageWrapper {
    return &MessageWrapper{msg: m}
}

func (m *MessageWrapper) InnerMessage() Message {
    return m.msg
}

func (m *MessageWrapper) Encode() []byte {
    bytes := m.msg.Encode()
    encoder := codec.Encoder{}
    bytes = append(bytes, encoder.EncodeTLV(m.GetValidationAlgorithm())...)
    bytes = append(bytes, encoder.EncodeTLV(m.GetValidationPayload())...)

    return bytes
}

func CreateFromTLV(tlv []codec.TLV) (*MessageWrapper, error) {
    var inner Message
    var validationAlgorithm validation.ValidationAlgorithm
    var validationPayload validation.ValidationPayload
    var err error

    for _, root := range tlv {
        switch (root.Type()) {
        case codec.T_INTEREST:
            inner, err = interest.CreateFromTLV(root)
            break
        case codec.T_OBJECT:
            inner, err = content.CreateFromTLV(root)
            break
        case codec.T_VALALG:
            validationAlgorithm, err = validation.CreateFromTLV(root)
            if err != nil {
                return nil, err
            }
            break
        case codec.T_VALPAYLOAD:
            validationPayload = validation.NewValidationPayload(root.Value())
            break
        default:
            fmt.Println("Invalid type:", string(root.Type()))
            return nil, messageError{"Unable to create a message from the top-level TLV type " + string(root.Type())}
        }
    }

    wrapper := MessageWrapper {
        msg: inner,
        validationPayload: validationPayload,
        validationAlgorithm: validationAlgorithm,
    }

    return &wrapper, nil
}

// Messages can compute the hashes of their protected regions and their complete packet formats.
func (m *MessageWrapper) HashProtectedRegion(hasher hash.Hash) []byte {
    bytes := m.msg.Encode()
    encoder := codec.Encoder{}
    bytes = append(bytes, encoder.EncodeTLV(m.GetValidationAlgorithm())...)
    hasher.Write(bytes)
    return hasher.Sum(nil)
}

func (m *MessageWrapper) ComputeMessageHash(hasher hash.Hash) []byte {
    bytes := m.Encode()
    hasher.Write(bytes)
    return hasher.Sum(nil)
}

func (m *MessageWrapper) Identifier() string {
    if m.msg.Name().Length() > 0 {
        return m.msg.Name().String()
    } else {
        hash := m.ComputeMessageHash(sha256.New())
        return string(hash)
    }
}

func (m *MessageWrapper) NamelessIdentifier() string {
    hash := m.ComputeMessageHash(sha256.New())
    return string(hash)
}

func (m *MessageWrapper) SetValidationAlgorithm(va validation.ValidationAlgorithm) {
    m.validationAlgorithm = va
}

func (m *MessageWrapper) SetValidationPayload(vp validation.ValidationPayload) {
    m.validationPayload = vp
}

func (m *MessageWrapper) GetValidationAlgorithm() validation.ValidationAlgorithm {
    return m.validationAlgorithm
}

func (m *MessageWrapper) GetValidationPayload() validation.ValidationPayload {
    return m.validationPayload
}

func (m *MessageWrapper) GetPacketType() uint16 {
    return m.msg.GetPacketType()
}

func (m *MessageWrapper) Name() *name.Name {
    return m.msg.Name()
}

func (m *MessageWrapper) Payload() *payload.Payload {
    return m.msg.Payload()
}

func (m *MessageWrapper) PayloadType() uint16 {
    return m.msg.PayloadType()
}
