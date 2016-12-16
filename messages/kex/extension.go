package kex

import "github.com/chris-wood/spud/codec"

type KEXExtension struct {
    ExtType uint16
    ExtValue []byte
}

// TLV functions
func (e KEXExtension) Type() uint16 {
    return e.ExtType
}

func (e KEXExtension) TypeString() string {
    return "KEXExtension"
}

func (e KEXExtension) Length() uint16 {
    return uint16(len(e.ExtValue))
}

func (e KEXExtension) Value() []byte  {
    return e.ExtValue
}

func (e KEXExtension) Children() []codec.TLV {
    return make([]codec.TLV, 0) // we have no children
}

func (e KEXExtension) String() string {
    return e.TypeString()
}
