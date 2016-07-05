package interest

import "fmt"
import "hash"
import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages/name"
import typedhash "github.com/chris-wood/spud/messages/hash"
import "github.com/chris-wood/spud/messages/link"
import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/messages/validation"

type Interest struct {
    name name.Name
    keyId typedhash.Hash
    contentId typedhash.Hash
    dataPayload payload.Payload

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

func CreateWithName(name name.Name) Interest {
    return Interest{name: name}
}

// func CreateWithNameAndPayload(name *name.Name, payload []byte) *Interest {
//     return &Interest{name: name, payload: payload}
// }

func CreateFromLink(link link.Link) Interest {
    return Interest{name: link.Name(), keyId: link.KeyID(), contentId: link.ContentID()}
}

func CreateFromTLV(tlvs codec.TLV) (Interest, error) {
    var interest Interest
    var interestName name.Name
    var err error

    for _, tlv := range(tlvs.Children()) {
        if tlv.Type() == codec.T_NAME {
            interestName, err = name.CreateFromTLV(tlv)
            if err != nil {
                return interest, err
            }
        } else if tlv.Type() == codec.T_KEYID_REST {
            // pass
        } else if tlv.Type() == codec.T_HASH_REST {
            // pass
        } else {
            fmt.Printf("Unable to parse interest TLV type: %d\n", tlv.Type())
        }
    }

    return Interest{name: interestName}, nil
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
        length += i.dataPayload.Length() + 4
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
        value = append(value, e.EncodeTLV(i.dataPayload)...)
    }

    return value
}

func (i Interest) Children() []codec.TLV  {
    children := []codec.TLV{i.name, i.keyId, i.contentId}
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

func (i Interest) ComputeMessageHash(hasher hash.Hash) []byte {
    return make([]byte, 0)
}

func (i Interest) Encode() []byte {
    encoder := codec.Encoder{}
    bytes := encoder.EncodeTLV(i)
    return bytes
}

func (i Interest) HashSensitiveRegion(hasher hash.Hash) []byte {
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

    return value
}

func (i Interest) IsRequest() bool {
    return true
}

func (i Interest) Payload() payload.Payload {
    return i.dataPayload
}

func (i Interest) SetValidationAlgorithm(va validation.ValidationAlgorithm) {
    i.validationAlgorithm = va
}
