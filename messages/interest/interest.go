package interest

import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/hash"
import "github.com/chris-wood/spud/messages/link"

type Interest struct {
    name *name.Name
    keyId *hash.Hash
    contentId *hash.Hash
    payload []byte // TODO: make a payload TLV wrapper
}

// Constructors

func CreateWithName(name *name.Name) *Interest {
    return &Interest{name: name}
}

// func CreateWithNameAndPayload(name *name.Name, payload []byte) *Interest {
//     return &Interest{name: name, payload: payload}
// }

func CreateFromLink(link *link.Link) *Interest {
    return &Interest{name: link.Name(), keyId: link.KeyID(), contentId: link.ContentID()}
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
    if i.keyId != nil {
        length += i.keyId.Length() + 4
    }
    if i.contentId != nil {
        length += i.contentId.Length() + 4
    }
    if i.payload != nil {
        length += uint16(len(i.payload)) + 4
    }
    return length
}

func (i Interest) Value() []byte  {

    e := codec.Encoder{}
    value := e.EncodeTLV(i.name)
    if i.keyId != nil {
        value = append(value, e.EncodeTLV(i.keyId)...)
    }
    if i.contentId != nil {
        value = append(value, e.EncodeTLV(i.contentId)...)
    }
    if i.payload != nil {
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
