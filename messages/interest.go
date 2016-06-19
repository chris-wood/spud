package messages

import "github.com/chris-wood/spud/messages/name"

type Interest struct {
    name *name.Name
    payload []byte
}

// Constructors

// TODO

// TLV functions

func (i Interest) Type() uint16 {
    return uint16(codec.T_INTEREST)
}

func (i Interest) TypeString() string {
    return "Interest"
}

func (i Interest) Length() uint16 {
    // length := uint16(0)
    // for _, ns := range(n.Segments) {
    //     length += ns.Length() + 4
    // }
    // return length
    return 0
}

func (i Interest) Value() []byte  {
    value := make([]byte, 0)

    // e := codec.Encoder{}
    // for _, segment := range(n.Segments) {
    //     value = append(value, e.Encode(segment)...)
    // }

    return value
}

// Message functions

func (i Interest) ComputeMessageHash() []byte {
    return make([]byte, 0)
}
