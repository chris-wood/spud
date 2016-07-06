package publickey

import "github.com/chris-wood/spud/codec"
import "fmt"

type PublicKey struct {
    bytes []byte
}

type keyError struct {
    prob string
}

func (e keyError) Error() string {
    return fmt.Sprintf("%s", e.prob)
}

// Constructors

func Create(bytes []byte) PublicKey {
    return PublicKey{bytes: bytes}
}

func CreateFromTLV(tlv codec.TLV) (PublicKey, error) {
    return PublicKey{bytes: tlv.Value()}, nil
}

// TLV functions

func (pk PublicKey) Type() uint16 {
    return codec.T_PUBLICKEY
}

func (pk PublicKey) TypeString() string {
    return "PublicKey"
}

func (pk PublicKey) Length() uint16 {
    return uint16(len(pk.bytes))
}

func (pk PublicKey) Value() []byte  {
    return pk.bytes
}

func (pk PublicKey) Children() []codec.TLV {
    return nil
}

// String functions

func (pk PublicKey) String() string {
    return string(pk.bytes)
}
