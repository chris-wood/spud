package flic

import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/hash"
import "github.com/chris-wood/spud/codec"

const FLICPointerLimit int = 10

type HashGroupMetadata struct {
    Locator name.Name
    OverallByteCount Size
    OverallDataDigest hash.Hash
}

type HashGroup struct {
    metadata *HashGroupMetadata
    Pointers []SizedPointer
}

func CreateEmptyHashGroup() *HashGroup {
    return &HashGroup{nil, make([]SizedPointer, 0)}
}

func (h *HashGroup) AddPointer(p SizedPointer) bool {
    if len(h.Pointers) < FLICPointerLimit {
        h.Pointers = append(h.Pointers, p)
        return true
    }
    return false
}

// TLV functions

func (g HashGroup) Type() uint16 {
    return uint16(codec.T_HASHGROUP)
}

func (g HashGroup) TypeString() string {
    return "HashGroup"
}

func (g HashGroup) Value() []byte {
    value := make([]byte, 0)

    if g.metadata != nil {
        // XXX
    }

    e := codec.Encoder{}
    for _, pointer := range(g.Pointers) {
        value = append(value, e.EncodeTLV(pointer)...)
    }

    return value
}

func (g HashGroup) Length() uint16 {
    e := codec.Encoder{}
    encodedValue := e.EncodeTLV(g)
    return uint16(len(encodedValue))
}

func (g HashGroup) Children() []codec.TLV {
    children := []codec.TLV{}

    if g.metadata != nil {
        // XXX
    }

    for _, pointer := range(g.Pointers) {
        children = append(children, pointer)
    }

    return children
}

func (g HashGroup) String() string {
    return "HashGroup"
}
