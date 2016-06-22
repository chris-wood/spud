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

func Create(length int, hash hash.Hash) *Hash {
    bytes := make([]byte, length)
    bytes = hash.Sum(bytes)

    // TODO: do some error checking on the length here?

    return &Hash{digest: bytes}
}

// TLV functions

func (h Hash) Type() uint16 {
    return uint16(codec.T_HASH)
}

func (h Hash) TypeString() string {
    return "Hash"
}

func (h Hash) Length() uint16 {
    return len(h.digest)
}

func (h Hash) Value() []byte  {
    return h.digest
}

func (h Hash) Children() []codec.TLVInterface {
    return nil
}

// String functions

func (h Hash) String() string {
    return string(h.digest)
}
