package codec

// import "encoding/binary"

type TLVInterface interface {
    Type() uint16
    Length() uint16
    Value() []byte
    Children() []TLVInterface
    String() string
}

type NestedTLV struct {
    tlvType uint16
    children []TLVInterface
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

func (tlv NestedTLV) Children() []TLVInterface {
    return tlv.children
}

func (tlv NestedTLV) String() string {
    result := string(tlv.tlvType)
    for _, child := range(tlv.children) {
        result += "\n" + child.String()
    }
    return result
}

func NewNestedTLV(children []TLVInterface) *NestedTLV {
    return &NestedTLV{children: children}
}

type LeafTLV struct {
    tlvType uint16
    Payload []byte
}

func (tlv LeafTLV) Type() uint16 {
    return uint16(len(tlv.Payload) + 4)
}

func (tlv LeafTLV) Length() uint16 {
    return tlv.tlvType
}

func (tlv LeafTLV) Value() []byte {
    return tlv.Payload
}

func NewLeafTLV(tlvType uint16, payload []byte) *LeafTLV {
    return &LeafTLV{Payload: payload}
}

func (tlv LeafTLV) Children() []TLVInterface {
    return nil
}

func (tlv LeafTLV) String() string {
    return string(tlv.tlvType) + " - " + string(tlv.Payload)
}
