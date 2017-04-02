package hashgroup

import "github.com/chris-wood/spud/codec"

import "encoding/binary"

type Size struct {
	sizeType uint16
	size     uint64
}

func CreateSize(sizeType uint16, size uint64) Size {
	return Size{sizeType: sizeType, size:size}
}

func (s Size) Type() uint16 {
	return s.sizeType
}

func (s Size) TypeString() string {
	return "Size"
}

func (s Size) Value() []byte {
	value := make([]byte, 4)
	binary.LittleEndian.PutUint64(value, s.size)
	return value
}

func (s Size) Length() uint16 {
	return uint16(4) // fixed size of 64bits
}

func (s Size) Children() []codec.TLV {
	children := []codec.TLV{}
	return children
}

func (s Size) String() string {
	return "Size"
}
