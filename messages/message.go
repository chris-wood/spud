package messages

import "github.com/chris-wood/spud/codec"

type Message interface {
    CreateFromTLV(tlv codec.TLV) *Message

    HashSensitiveRegion() []byte
    ComputeMessageHash() []byte
    Encode() []byte

    // XX: this should take a signer as input
    // TagAndEncode() []byte
}

func CreateFromTLV(tlv []codec.TLV) Message {
    var result Message

    return result
}
