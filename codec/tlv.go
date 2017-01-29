package codec

// import "fmt"

import "encoding/json"

type TLV interface {
    Type() uint16
    Length() uint16
    Value() []byte

    Children() []TLV
    String() string
}

type NestedTLV struct {
    tlvType uint16 `json:"type"`
    children []TLV `json:"children"`
}

func (tlv NestedTLV) Type() uint16 {
    return tlv.tlvType
}

func (tlv NestedTLV) Length() uint16 {
    tlvLength := uint16(0)
    for _, child := range(tlv.children) {
        tlvLength += 4 + child.Length()
    }
    return tlvLength
}

func (tlv NestedTLV) Value() []byte {
    e := Encoder{}
    childrenBytes := e.Encode(tlv.children)
    return childrenBytes
}

func (tlv NestedTLV) Children() []TLV {
    return tlv.children
}

func (tlv NestedTLV) String() string {
    result, err := json.Marshal(tlv)
    if err == nil {
        return string(result)
    }
    return err.Error()
}

// func (tlv NestedTLV) DisplayTLV(indents int) {
//     for i := 0; i < indents; i++ {
//         fmt.Printf("  ")
//     }
//     fmt.Printf("%d %d\n", tlv.Type(), tlv.Length())
//     for _, c := range(tlv.children) {
//         c.DisplayTLV(indents + 1)
//     }
// }

func NewNestedTLV(tlvType uint16, children []TLV) NestedTLV {
    return NestedTLV{tlvType: tlvType, children: children}
}

type LeafTLV struct {
    tlvType uint16 `json:"type"`
    Payload []byte `json:"payload"`
}

func (tlv LeafTLV) Type() uint16 {
    return tlv.tlvType
}

func (tlv LeafTLV) Length() uint16 {
    return uint16(len(tlv.Payload))
}

func (tlv LeafTLV) Value() []byte {
    return tlv.Payload
}

func NewLeafTLV(tlvType uint16, payload []byte) LeafTLV {
    return LeafTLV{tlvType: tlvType, Payload: payload}
}

func (tlv LeafTLV) Children() []TLV {
    return nil
}

func (tlv LeafTLV) String() string {
    result, err := json.Marshal(tlv)
    if err == nil {
        return string(result)
    }
    return err.Error()
}

// func (tlv LeafTLV) DisplayTLV(indents int) {
//     for i := 0; i < indents; i++ {
//         fmt.Printf("  ")
//     }
//     fmt.Printf("%d %d: %s\n", tlv.Type(), tlv.Length(), string(tlv.Payload))
// }
