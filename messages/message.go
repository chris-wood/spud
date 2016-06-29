package messages

import "github.com/chris-wood/spud/codec"

type Message interface {
    Identifier() string
    HashSensitiveRegion() []byte
    ComputeMessageHash() []byte
    Encode() []byte

    // XXX: this should take a signer as input
    // TagAndEncode() []byte
}

// XXX: create the right type of message here
func CreateFromTLV(tlv []codec.TLV) (Message, error) {
    var result Message

    return result, nil
}
