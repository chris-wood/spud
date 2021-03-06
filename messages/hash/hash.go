package hash

import "fmt"

import (
	"github.com/chris-wood/spud/codec"
	"crypto/subtle"
)

type HashType uint16

const HashTypeSHA256 HashType = 0x0001
const HashTypeSHA512 HashType = 0x0002

type Hash struct {
	hashType HashType
	digest   []byte
}

type hashError struct {
	prob string
}

func (e hashError) Error() string {
	return fmt.Sprintf("%s", e.prob)
}

// Constructors

func Create(hashType HashType, value []byte) Hash {
	return Hash{hashType, value}
}

func CreateFromTLV(tlv codec.TLV) (Hash, error) {
	var result Hash
	return result, nil
}

// TLV functions

func (h Hash) Type() uint16 {
	return uint16(codec.T_HASH)
}

func (h Hash) TypeString() string {
	return "Hash"
}

func (h Hash) Length() uint16 {
	return uint16(len(h.digest))
}

func (h Hash) Value() []byte {
	return h.digest
}

func (h Hash) Children() []codec.TLV {
	return nil
}

func (h Hash) String() string {
	return string(h.digest)
}

// Utility functions

func (h Hash) Equals(other Hash) bool {
	if h.hashType == other.hashType {
		if subtle.ConstantTimeCompare(h.digest, other.digest) == 0 {
			return true
		}
	}
	return false
}