package content

import "fmt"
import "hash"
import "crypto/sha256"

import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/kex"
import "github.com/chris-wood/spud/messages/validation"
import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/codec"

type Content struct {
    name name.Name
    dataPayload payload.Payload
    payloadType uint8

    containers []codec.TLV

    // Validation information
    validationAlgorithm validation.ValidationAlgorithm
    validationPayload validation.ValidationPayload
}

type contentError struct {
    prob string
}

func (e contentError) Error() string {
    return fmt.Sprintf("%s", e.prob)
}

// Constructors

func CreateWithName(stubName name.Name) *Content {
    var dataPayload payload.Payload
    return &Content{name: stubName, dataPayload: dataPayload, payloadType: codec.T_PAYLOADTYPE_DATA, containers: make([]codec.TLV, 0)}
}

func CreateWithPayload(dataPayload payload.Payload) *Content {
    var name name.Name
    return &Content{name: name, dataPayload: dataPayload, payloadType: codec.T_PAYLOADTYPE_DATA, containers: make([]codec.TLV, 0)}
}

func CreateWithNameAndPayload(name name.Name, dataPayload payload.Payload) *Content {
    return &Content{name: name, dataPayload: dataPayload, payloadType: codec.T_PAYLOADTYPE_DATA, containers: make([]codec.TLV, 0)}
}

func CreateWithNameAndTypedPayload(name name.Name, payloadType uint8, dataPayload payload.Payload) *Content {
    return &Content{name: name, dataPayload: dataPayload, payloadType: payloadType, containers: make([]codec.TLV, 0)}
}

// func CreateWithNameAndLink(name, *name.Name, payload []byte) *Content {
//     return &Content{name: name, payload: payload}
// }
//
// func CreateWithNameAndKey(name, *name.Name, payload []byte) *Content {
//     return &Content{name: name, payload: payload}
// }

func CreateFromTLV(topLevelTLV codec.TLV) (*Content, error) {
    var result Content
    var contentName name.Name
    var dataPayload payload.Payload
    var validationAlgorithm validation.ValidationAlgorithm
    var validationPayload validation.ValidationPayload
    var err error

    containers := make([]codec.TLV, 0)

    for _, tlv := range(topLevelTLV.Children()) {
        if tlv.Type() == codec.T_NAME {
            contentName, err = name.CreateFromTLV(tlv)
            if err != nil {
                return &result, err
            }
        } else if tlv.Type() == codec.T_PAYLOAD {
            dataPayload = payload.Create(tlv.Value())
        } else if tlv.Type() == codec.T_KEX {
            kex, err := kex.CreateFromTLV(tlv)
            if err != nil {
                return nil, contentError{"Unable to decode the KEX TLV"}
            }
            containers = append(containers, kex)
        } else if tlv.Type() == codec.T_VALALG {
            validationAlgorithm, err = validation.CreateFromTLV(tlv)
            if err != nil {
                return &result, err
            }
        } else if tlv.Type() == codec.T_VALPAYLOAD {
            validationPayload = validation.NewValidationPayload(tlv.Value())
        } else {
            fmt.Printf("Unable to parse content TLV type: %d\n", tlv.Type())
        }
    }

    return &Content{name: contentName, dataPayload: dataPayload, validationAlgorithm: validationAlgorithm,
        validationPayload: validationPayload, containers: containers}, nil
}

// Containers

func (c *Content) AddContainer(container codec.TLV) {
    c.containers = append(c.containers, container)
}

func (c *Content) GetContainer(containerType uint16) (codec.TLV, error) {
    var container codec.TLV
    for _, test := range(c.containers) {
        if test.Type() == containerType {
            return test, nil
        }
    }
    return container, contentError{"No such container"}
}

// TLV functions

func (c Content) Type() uint16 {
    return uint16(codec.T_OBJECT)
}

func (c Content) TypeString() string {
    return "Content"
}

func (c Content) Length() uint16 {
    length := uint16(0)

    if c.name.Length() > 0 {
        length += c.name.Length() + 4
    }

    if c.dataPayload.Length() > 0 {
        length += c.dataPayload.Length() + 4
    }

    if c.validationAlgorithm.Length() > 0 {
        length += c.validationAlgorithm.Length() + 4
    }

    if c.validationPayload.Length() > 0 {
        length += c.validationPayload.Length() + 4
    }

    for _, container := range(c.containers) {
        length += container.Length() + 4
    }

    return length
}

func (c Content) Value() []byte {
    e := codec.Encoder{}
    value := make([]byte, 0)

    if c.name.Length() > 0 {
        value = append(value, e.EncodeTLV(c.name)...)
    }

    if c.dataPayload.Length() > 0 {
        value = append(value, e.EncodeTLV(c.dataPayload)...)
    }

    if c.validationAlgorithm.Length() > 0 {
        value = append(value, e.EncodeTLV(c.validationAlgorithm)...)
    }

    if c.validationPayload.Length() > 0 {
        value = append(value, e.EncodeTLV(c.validationPayload)...)
    }

    for _, container := range(c.containers) {
        value = append(value, e.EncodeTLV(container)...)
    }

    return value
}

func (c Content) Children() []codec.TLV {
    children := []codec.TLV{c.name, c.dataPayload, c.validationAlgorithm, c.validationPayload}
    for _, container := range(c.containers) {
        children = append(children, container)
    }
    return children
}

func (c Content) String() string {
    return c.Identifier()
}

// Message functions

func (c Content) ComputeMessageHash(hasher hash.Hash) []byte {
    encoder := codec.Encoder{}
    bytes := encoder.EncodeTLV(c)
    hasher.Write(bytes)
    return hasher.Sum(nil)
}

func (c Content) Encode() []byte {
    encoder := codec.Encoder{}
    bytes := encoder.EncodeTLV(c)
    return bytes
}

func (c Content) Name() name.Name {
    return c.name
}

func (c Content) Identifier() string {
    if c.name.Length() > 0 {
        return c.name.String()
    } else {
        hash := c.ComputeMessageHash(sha256.New())
        return string(hash)
    }
}

func (c Content) NamelessIdentifier() string {
    hash := c.ComputeMessageHash(sha256.New())
    return string(hash)
}

func (c Content) HashProtectedRegion(hasher hash.Hash) []byte {
    encoder := codec.Encoder{}

    value := make([]byte, 0)
    if c.name.Length() > 0 {
        value = append(value, encoder.EncodeTLV(c.name)...)
    }
    if c.dataPayload.Length() > 0 {
        value = append(value, encoder.EncodeTLV(c.dataPayload)...)
    }
    if c.validationAlgorithm.Length() > 0 {
        value = append(value, encoder.EncodeTLV(c.validationAlgorithm)...)
    }

    hasher.Write(value)
    return hasher.Sum(nil)
}

func (c Content) GetPacketType() uint16 {
    return codec.T_OBJECT
}

func (c Content) Payload() payload.Payload {
    return c.dataPayload
}

func (c Content) PayloadType() uint8 {
    return c.payloadType
}

func (c *Content) SetValidationAlgorithm(va validation.ValidationAlgorithm) {
    c.validationAlgorithm = va
}

func (c *Content) SetValidationPayload(vp validation.ValidationPayload) {
    c.validationPayload = vp
}

func (c Content) GetValidationAlgorithm() validation.ValidationAlgorithm {
    return c.validationAlgorithm
}

func (c Content) GetValidationPayload() validation.ValidationPayload {
    return c.validationPayload
}
