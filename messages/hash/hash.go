package hash

import "github.com/chris-wood/spud/codec"
import "hash"
import "fmt"

// import "encoding/json"

type Hash struct {
    digest []byte
}

type hashError struct {
    prob string
}

func (e hashError) Error() string {
    return fmt.Sprintf("%s", e.prob)
}

// Constructors

func Create(hash hash.Hash) Hash {
    return Hash{digest: hash.Sum(nil)}
}

func CreateTruncated(hash hash.Hash, truncatedLength int) Hash {
    bytes := make([]byte, truncatedLength)
    bytes = hash.Sum(bytes)
    return Hash{digest: bytes}
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

func (h Hash) Value() []byte  {
    return h.digest
}

func (h Hash) Children() []codec.TLV {
    return nil
}

func (h Hash) String() string {
    return string(h.digest)
}
