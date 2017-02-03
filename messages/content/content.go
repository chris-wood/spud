package content

import "fmt"

import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/kex"
import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/codec"

type Content struct {
    name *name.Name
    dataPayload payload.Payload
    payloadType uint16

    containers []codec.TLV
}

type contentError struct {
    prob string
}

func (e contentError) Error() string {
    return fmt.Sprintf("%s", e.prob)
}

// Constructors

func CreateWithName(stubName *name.Name) *Content {
    var dataPayload payload.Payload
    return &Content{name: stubName, dataPayload: dataPayload, payloadType: codec.T_PAYLOADTYPE_DATA, containers: make([]codec.TLV, 0)}
}

func CreateWithPayload(dataPayload payload.Payload) *Content {
    return &Content{name: nil, dataPayload: dataPayload, payloadType: codec.T_PAYLOADTYPE_DATA, containers: make([]codec.TLV, 0)}
}

func CreateWithNameAndPayload(name *name.Name, dataPayload payload.Payload) *Content {
    return &Content{name: name, dataPayload: dataPayload, payloadType: codec.T_PAYLOADTYPE_DATA, containers: make([]codec.TLV, 0)}
}

func CreateWithNameAndTypedPayload(name *name.Name, payloadType uint16, dataPayload payload.Payload) *Content {
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
    var contentName *name.Name
    var dataPayload payload.Payload
    var err error

    containers := make([]codec.TLV, 0)

    for _, tlv := range(topLevelTLV.Children()) {
        if tlv.Type() == codec.T_NAME {
            contentName, err = name.CreateFromTLV(tlv)
            if err != nil {
                return nil, err
            }
        } else if tlv.Type() == codec.T_PAYLOAD {
            dataPayload = payload.Create(tlv.Value())
        } else if tlv.Type() == codec.T_KEX {
            kex, err := kex.CreateFromTLV(tlv)
            if err != nil {
                return nil, contentError{"Unable to decode the KEX TLV"}
            }
            containers = append(containers, kex)
        } else {
            fmt.Printf("Unable to parse content TLV type: %d\n", tlv.Type())
        }
    }

    return &Content{name: contentName, dataPayload: dataPayload, containers: containers}, nil
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

    if c.name != nil && c.name.Length() > 0 {
        length += c.name.Length() + 4
    }

    if c.dataPayload.Length() > 0 {
        length += c.dataPayload.Length() + 4
    }

    for _, container := range(c.containers) {
        length += container.Length() + 4
    }

    return length
}

func (c Content) Value() []byte {
    e := codec.Encoder{}
    value := make([]byte, 0)

    if c.name != nil && c.name.Length() > 0 {
        value = append(value, e.EncodeTLV(c.name)...)
    }

    if c.dataPayload.Length() > 0 {
        value = append(value, e.EncodeTLV(c.dataPayload)...)
    }

    for _, container := range(c.containers) {
        value = append(value, e.EncodeTLV(container)...)
    }

    return value
}

func (c Content) Children() []codec.TLV {
    children := []codec.TLV{c.name, c.dataPayload}
    for _, container := range(c.containers) {
        children = append(children, container)
    }
    return children
}

func (c Content) String() string {
    // return Identifier(c)
    return c.name.String()
}

// Message functions

func (c Content) Encode() []byte {
    encoder := codec.Encoder{}
    bytes := encoder.EncodeTLV(c)
    return bytes
}

func (c Content) Name() *name.Name {
    return c.name
}

func (c Content) GetPacketType() uint16 {
    return codec.T_OBJECT
}

func (c Content) Payload() *payload.Payload {
    return &c.dataPayload
}

func (c Content) PayloadType() uint16 {
    return c.payloadType
}
