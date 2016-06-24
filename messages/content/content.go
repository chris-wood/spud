package content

import "github.com/chris-wood/spud/messages/name"

type Content struct {
    name name.Name

    // TODO: this must be a "payload interface"
    payload []byte

    // TODO: include the validation stuff
}

// Constructors

func CreateWithPayload(payload []byte) Content {
    return Content{name: nil, payload: payload}
}

func CreateWithNameAndPayload(name, *name.Name, payload []byte) Content {
    return Content{name: name, payload: payload}
}

// func CreateWithNameAndLink(name, *name.Name, payload []byte) *Content {
//     return &Content{name: name, payload: payload}
// }
//
// func CreateWithNameAndKey(name, *name.Name, payload []byte) *Content {
//     return &Content{name: name, payload: payload}
// }

// TLV functions

func (c Content) Type() uint16 {
    return uint16(codec.T_CONTENT)
}

func (c Content) TypeString() string {
    return "Content"
}

func (c Content) Length() uint16 {
    // length := uint16(0)
    // for _, ns := range(n.Segments) {
    //     length += ns.Length() + 4
    // }
    // return length
    return 0
}

func (c Content) Value() []byte  {
    value := make([]byte, 0)

    // e := codec.Encoder{}
    // for _, segment := range(n.Segments) {
    //     value = append(value, e.Encode(segment)...)
    // }

    return value
}

// Message functions

func (c Content) ComputeMessageHash() []byte {
    return make([]byte, 0)
}
