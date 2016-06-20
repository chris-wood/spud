package codec

// import "encoding/binary"

type TLVInterface interface {
    Type() uint16
    Length() uint16
    Value() []byte
}

type NestedTLV struct {
    tlvType uint16
    Children []TLVInterface
}

func (tlv NestedTLV) Type() uint16 {
    return tlv.tlvType
}

func (tlv NestedTLV) Length() uint16 {
    tlvLength := uint16(0)
    for _, child := range(tlv.Children) {
        tlvLength += 4 + child.Length()
    }
    return tlvLength
}

func (tlv NestedTLV) Value() []byte {
    e := Encoder{}
    childrenBytes := e.Encode(tlv.Children)
    return childrenBytes
}

func NewNestedTLV(children []TLVInterface) *NestedTLV {
    return &NestedTLV{Children: children}
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
