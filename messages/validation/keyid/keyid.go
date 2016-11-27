package keyid

import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages/hash"
import "fmt"

type KeyId struct {
    keyDigest hash.Hash
}

type keyError struct {
    prob string
}

func (e keyError) Error() string {
    return fmt.Sprintf("%s", e.prob)
}

// Constructors

func Create(keyDigest hash.Hash) KeyId {
    return KeyId{keyDigest: keyDigest}
}

func CreateFromTLV(tlv codec.TLV) (KeyId, error) {
    var result KeyId

    digest, err := hash.CreateFromTLV(tlv.Children()[0])
    if err != nil {
        return result, err
    }

    return KeyId{keyDigest: digest}, nil
}

// TLV functions

func (kid KeyId) Type() uint16 {
    return codec.T_KEYID
}

func (kid KeyId) TypeString() string {
    return "KeyId"
}

func (kid KeyId) Length() uint16 {
    return kid.keyDigest.Length() + 4
}

func (kid KeyId) Value() []byte {
    e := codec.Encoder{}
    return e.EncodeTLV(kid.keyDigest)
}

func (kid KeyId) Children() []codec.TLV {
    return []codec.TLV{kid.keyDigest}
}

func (kid KeyId) String() string {
    return kid.keyDigest.String()
}
