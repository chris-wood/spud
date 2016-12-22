package interest

import "fmt"
import "hash"
import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/kex"
import typedhash "github.com/chris-wood/spud/messages/hash"
import "github.com/chris-wood/spud/messages/link"
import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/messages/validation"

type Interest struct {
    name name.Name
    keyId typedhash.Hash
    contentId typedhash.Hash

    // Payload and its type
    payloadType uint8
    dataPayload payload.Payload

    // codec.TLVs
    containers []codec.TLV

    // KEX signalling and encryption information
    kexMessage kex.KEX

    // Validation information
    validationAlgorithm validation.ValidationAlgorithm
    validationPayload validation.ValidationPayload
}

type interestError struct {
    prob string
}

func (e interestError) Error() string {
    return fmt.Sprintf("%s", e.prob)
}

// Constructors

func CreateWithName(name name.Name) *Interest {
    return &Interest{name: name, containers: make([]codec.TLV, 0)}
}

func CreateWithNameAndPayload(name name.Name, payloadType uint8, payloadValue payload.Payload) *Interest {
    return &Interest{name: name, payloadType: payloadType, dataPayload: payloadValue}
}

func CreateFromLink(link link.Link) *Interest {
    return &Interest{name: link.Name(), keyId: link.KeyID(), contentId: link.ContentID(), containers: make([]codec.TLV, 0)}
}

func CreateFromTLV(tlvs codec.TLV) (*Interest, error) {
    var interest Interest
    var interestName name.Name
    var err error

    containers := make([]codec.TLV, 0)

    for _, tlv := range(tlvs.Children()) {
        if tlv.Type() == codec.T_NAME {
            interestName, err = name.CreateFromTLV(tlv)
            if err != nil {
                return &interest, err
            }
        } else if tlv.Type() == codec.T_KEX {
            kex, err := kex.CreateFromTLV(tlv)
            if err != nil {
                return nil, interestError{"Unable to decode the KEX TLV"}
            }
            containers = append(containers, kex)
        } else if tlv.Type() == codec.T_PAYLDTYPE {
            // pass
        } else if tlv.Type() == codec.T_PAYLOAD {
            // pass
        } else if tlv.Type() == codec.T_KEYID_REST {
            // pass
        } else if tlv.Type() == codec.T_HASH_REST {
            // pass
        } else {
            fmt.Printf("Unable to parse interest TLV type: %d\n", tlv.Type())
        }
    }

    return &Interest{name: interestName, containers: containers}, nil
}

// codec.TLVs

func (i *Interest) AddContainer(container codec.TLV) {
    i.containers = append(i.containers, container)
}

func (i *Interest) GetContainer(containerType uint16) (codec.TLV, error) {
    var container codec.TLV
    for _, test := range(i.containers) {
        if test.Type() == containerType {
            return test, nil
        }
    }
    return container, interestError{"No such container"}
}

// TLV functions

func (i Interest) Type() uint16 {
    return uint16(codec.T_INTEREST)
}

func (i Interest) TypeString() string {
    return "Interest"
}

func (i Interest) Length() uint16 {
    length := i.name.Length() + 4

    if i.keyId.Length() > 0 {
        length += i.keyId.Length() + 4
    }
    if i.contentId.Length() > 0 {
        length += i.contentId.Length() + 4
    }
    if i.dataPayload.Length() > 0 {
        // length += 5 // for payload type and TLV header
        length += i.dataPayload.Length() + 4
    }
    for _, container := range(i.containers) {
        length += container.Length() + 4
    }

    return length
}

func (i Interest) Value() []byte  {
    e := codec.Encoder{}
    value := e.EncodeTLV(i.name)

    if i.keyId.Length() > 0 {
        value = append(value, e.EncodeTLV(i.keyId)...)
    }
    if i.contentId.Length() > 0 {
        value = append(value, e.EncodeTLV(i.contentId)...)
    }
    if i.dataPayload.Length() > 0 {
        // XXX: append the payload type, encoded, here
        value = append(value, e.EncodeTLV(i.dataPayload)...)
    }
    for _, container := range(i.containers) {
        value = append(value, e.EncodeTLV(container)...)
    }

    return value
}

func (i Interest) Children() []codec.TLV  {
    children := []codec.TLV{i.name, i.keyId, i.contentId}
    for _, container := range(i.containers) {
        children = append(children, container)
    }
    return children
}

func (i Interest) String() string  {
    return i.name.String()
}

// Message functions

func (i Interest) Name() name.Name {
    return i.name
}

func (i Interest) Identifier() string {
    return i.name.String()
}

func (i Interest) NamelessIdentifier() string {
    result := ""
    if i.keyId.Length() > 0 {
        result += string(i.keyId.String())
    }
    if i.contentId.Length() > 0 {
        result += string(i.contentId.String())
    }
    return result
}

func (i Interest) ComputeMessageHash(hasher hash.Hash) []byte {
    return make([]byte, 0)
}

func (i Interest) Encode() []byte {
    encoder := codec.Encoder{}
    bytes := encoder.EncodeTLV(i)
    return bytes
}

func (i Interest) HashProtectedRegion(hasher hash.Hash) []byte {
    encoder := codec.Encoder{}

    value := encoder.EncodeTLV(i.name)
    if i.keyId.Length() > 0 {
        value = append(value, encoder.EncodeTLV(i.keyId)...)
    }
    if i.contentId.Length() > 0 {
        value = append(value, encoder.EncodeTLV(i.contentId)...)
    }
    if i.dataPayload.Length() > 0 {
        value = append(value, encoder.EncodeTLV(i.dataPayload)...)
    }
    if i.validationAlgorithm.Length() > 0 {
        value = append(value, encoder.EncodeTLV(i.validationAlgorithm)...)
    }

    hasher.Write(value)
    return hasher.Sum(nil)
}

func (i Interest) GetPacketType() uint16 {
    return codec.T_INTEREST
}

func (i Interest) Payload() payload.Payload {
    return i.dataPayload
}

func (i Interest) PayloadType() uint8 {
    return i.payloadType
}

func (i *Interest) SetValidationAlgorithm(va validation.ValidationAlgorithm) {
    i.validationAlgorithm = va
}

func (i *Interest) SetValidationPayload(vp validation.ValidationPayload) {
    i.validationPayload = vp
}

func (i Interest) GetValidationAlgorithm() validation.ValidationAlgorithm {
    return i.validationAlgorithm
}

func (i Interest) GetValidationPayload() validation.ValidationPayload {
    return i.validationPayload
}
