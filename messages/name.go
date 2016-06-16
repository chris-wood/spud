package messages

import "github.com/chris-wood/spud/codec"
// import "encoding/json"

type Name struct {
    Segments []NameSegment `json:"segments"`
}

func (n Name) Type() uint16 {
    return uint16(1)
}

func (n Name) Length() uint16 {
    length := uint16(0)
    for _, ns := range(n.Segments) {
        length += ns.Length() + 4
    }
    return length
}

func (n Name) Value() []byte  {
    value := make([]byte, 0)

    e := codec.Encoder{}
    for _, segment := range(n.Segments) {
        value = append(value, e.Encode(segment)...)
    }

    return value
}
