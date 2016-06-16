package codec

// import "encoding/binary"
// import "github.com/chris-wood/spud/codec"

type TLV interface {
    Type() uint16
    Length() uint16
    Value() []byte
    // ToJSON() string
}

type NestedTLV struct {
    tlvType uint16
    Children []TLV
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
    children := make([]byte, 0)
    e := Encoder{}
    for _, child := range(tlv.Children) {
        children = append(children, e.Encode(child)...)
    }
    return children
}

func NewNestedTLV(children []TLV) *NestedTLV {
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

func NewLeafTLV(payload []byte) *LeafTLV {
    return &LeafTLV{Payload: payload}
}
