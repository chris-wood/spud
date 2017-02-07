package flic

import "fmt"

import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/codec"

type FLIC struct {
    rootName *name.Name
    groups []*HashGroup
}

type flicError struct {
    prob string
}

func (e flicError) Error() string {
    return fmt.Sprintf("%s", e.prob)
}

func CreateFLICFromHashGroup(group *HashGroup) *FLIC {
    return &FLIC{rootName: nil, groups: []*HashGroup{group}}
}

// Constructors

func CreateFromTLV(topLevelTLV codec.TLV) (*FLIC, error) {
    // var result FLIC
    // var err error

    // containers := make([]codec.TLV, 0)
    // for _, tlv := range(topLevelTLV.Children()) {
        // if tlv.Type() == codec.T_NAME {
        //     contentName, err = name.CreateFromTLV(tlv)
        //     if err != nil {
        //         return &result, err
        //     }
        // } else if tlv.Type() == codec.T_PAYLOAD {
        //     dataPayload = payload.Create(tlv.Value())
        // } else if tlv.Type() == codec.T_KEX {
        //     kex, err := kex.CreateFromTLV(tlv)
        //     if err != nil {
        //         return nil, contentError{"Unable to decode the KEX TLV"}
        //     }
        //     containers = append(containers, kex)
        // } else {
        //     fmt.Printf("Unable to parse content TLV type: %d\n", tlv.Type())
        // }
    // }

    return &FLIC{rootName: nil, groups: make([]*HashGroup, 0)}, nil
}

// Containers

func (f *FLIC) AddContainer(container codec.TLV) {
    // no-op
}

func (f *FLIC) GetContainer(containerType uint16) (codec.TLV, error) {
    var container codec.TLV
    return container, flicError{"FLIC types do not support containers"}
}

// TLV functions

func (f FLIC) Type() uint16 {
    return uint16(codec.T_MANIFEST)
}

func (f FLIC) TypeString() string {
    return "FLIC"
}

func (f FLIC) Value() []byte {
    value := make([]byte, 0)

    e := codec.Encoder{}
    if f.rootName != nil {
        value = append(value, e.EncodeTLV(f.rootName)...)
    }

    for _, group := range(f.groups) {
        value = append(value, e.EncodeTLV(group)...)
    }

    return value
}

func (f FLIC) Length() uint16 {
    e := codec.Encoder{}
    encodedValue := e.EncodeTLV(f)
    return uint16(len(encodedValue))
}

func (f FLIC) Children() []codec.TLV {
    children := []codec.TLV{}

    if f.rootName != nil {
        children = append(children, f.rootName)
    }
    for _, group := range(f.groups) {
        children = append(children, group)
    }

    return children
}

func (f FLIC) String() string {
    return "FLIC"
}

// Message functions

func (f FLIC) Encode() []byte {
    encoder := codec.Encoder{}
    bytes := encoder.EncodeTLV(f)
    return bytes
}

func (f FLIC) Name() *name.Name {
    return f.rootName
}

func (f FLIC) GetPacketType() uint16 {
    return codec.T_OBJECT
}

func (f FLIC) Payload() *payload.Payload {
    return nil
}

func (f FLIC) PayloadType() uint16 {
    return codec.T_PAYLOADTYPE_MANIFEST
}
