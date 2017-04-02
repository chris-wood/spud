package hashgroup

import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages/hash"

const hashGroupDataPointerType uint16 = codec.T_DATA_POINTER
const hashGroupManifestPointerType uint16 = codec.T_MANIFEST_POINTER

type SizedDataPointer struct {
	size    Size
	ptrHash hash.Hash
}

func CreateSizedDataPointer(size Size, ptrhash hash.Hash) *SizedDataPointer {
	return &SizedDataPointer{size: size, ptrHash:ptrhash}
}

func (p SizedDataPointer) GetPointerType() uint16 {
	return hashGroupDataPointerType
}

func (p SizedDataPointer) GetSize() Size {
	return p.size
}

func (p SizedDataPointer) GetPointer() hash.Hash {
	return p.ptrHash
}

func (p SizedDataPointer) Type() uint16 {
	return p.GetPointerType()
}

func (p SizedDataPointer) TypeString() string {
	return "DataPointer"
}

func (p SizedDataPointer) Value() []byte {
	value := make([]byte, 0)
	e := codec.Encoder{}
	value = append(value, e.EncodeTLV(p.GetSize())...)
	value = append(value, e.EncodeTLV(p.GetPointer())...)
	return value
}

func (p SizedDataPointer) Length() uint16 {
	e := codec.Encoder{}
	encodedValue := e.EncodeTLV(p)
	return uint16(len(encodedValue))
}

func (p SizedDataPointer) Children() []codec.TLV {
	children := []codec.TLV{p.GetSize(), p.GetPointer()}
	return children
}

func (p SizedDataPointer) String() string {
	return "DataPointer"
}

type SizedManifestPointer struct {
	size    Size
	ptrHash hash.Hash
}

func CreateSizedManifestPointer(size Size, ptrhash hash.Hash) *SizedManifestPointer {
	return &SizedManifestPointer{size: size, ptrHash:ptrhash}
}

func (p SizedManifestPointer) GetPointerType() uint16 {
	return hashGroupManifestPointerType
}

func (p SizedManifestPointer) GetSize() Size {
	return p.size
}

func (p SizedManifestPointer) GetPointer() hash.Hash {
	return p.ptrHash
}

func (p SizedManifestPointer) Type() uint16 {
	return p.GetPointerType()
}

func (p SizedManifestPointer) TypeString() string {
	return "ManifestPointer"
}

func (p SizedManifestPointer) Value() []byte {
	value := make([]byte, 0)
	e := codec.Encoder{}
	value = append(value, e.EncodeTLV(p.GetSize())...)
	value = append(value, e.EncodeTLV(p.GetPointer())...)
	return value
}

func (p SizedManifestPointer) Length() uint16 {
	e := codec.Encoder{}
	encodedValue := e.EncodeTLV(p)
	return uint16(len(encodedValue))
}

func (p SizedManifestPointer) Children() []codec.TLV {
	children := []codec.TLV{p.GetSize(), p.GetPointer()}
	return children
}

func (p SizedManifestPointer) String() string {
	return "ManifestPointer"
}

type SizedPointer interface {
	GetPointerType() uint16
	GetSize() Size
	GetPointer() hash.Hash

	// TLV functions
	Type() uint16
	TypeString() string
	Value() []byte
	Length() uint16
	Children() []codec.TLV
	String() string
}
