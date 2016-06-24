package interest

import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/hash"
import "github.com/chris-wood/spud/messages/link"

type Interest struct {
    name name.Name
    keyId hash.Hash
    contentId hash.Hash

    payload []byte // TODO: make a payload TLV wrapper

    // TODO: include the validation fields
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
    if len(i.payload) > 0 {
        length += uint16(len(i.payload)) + 4
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
    if len(i.payload) > 0 {
        value = append(value, i.payload...)
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

func (i Interest) ComputeMessageHash() []byte {
    return make([]byte, 0)
}

func (i Interest) Encode() []byte {
    encoder := codec.Encoder{}
    bytes := encoder.EncodeTLV(i)
    return bytes
}

func (i Interest) HashSensitiveRegion() []byte {
    encoder := codec.Encoder{}

    value := encoder.EncodeTLV(i.name)
    if i.keyId.Length() > 0 {
        value = append(value, encoder.EncodeTLV(i.keyId)...)
    }
    if i.contentId.Length() > 0 {
        value = append(value, encoder.EncodeTLV(i.contentId)...)
    }
    if len(i.payload) > 0 {
        value = append(value, i.payload...)
    }

    return value
}
