package content

import "fmt"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/codec"

type Content struct {
    name name.Name

    // TODO: this must be a "payload interface"
    payload []byte

    // TODO: include the validation stuff
}

type contentError struct {
    prob string
}

func (e contentError) Error() string {
    return fmt.Sprintf("%s", e.prob)
}

// Constructors

func CreateWithPayload(payload []byte) Content {
    var name name.Name
    return Content{name: name, payload: payload}
}

func CreateWithNameAndPayload(name name.Name, payload []byte) Content {
    return Content{name: name, payload: payload}
}

func CreateFromTLV(tlv []codec.TLV) (Content, error) {
    var result Content
    return result, contentError{"couldn't parse the content TLV"}
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
    return uint16(codec.T_OBJECT)
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

func (c Content) Children() []codec.TLV  {
    // XXX: need to include the payload here
    children := []codec.TLV{c.name}
    return children
}

func (c Content) String() string  {
    return c.Identifier()
}

// Message functions

func (c Content) ComputeMessageHash() []byte {
    return make([]byte, 0)
}

func (c Content) Encode() []byte {
    encoder := codec.Encoder{}
    bytes := encoder.EncodeTLV(c)
    return bytes
}

func (c Content) Identifier() string {
    return "TODO"
}

func (c Content) HashSensitiveRegion() []byte {
    return nil
}
