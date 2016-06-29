package messages

import "github.com/chris-wood/spud/codec"

type Message interface {
    Identifier() string
    HashSensitiveRegion() []byte
    ComputeMessageHash() []byte
    Encode() []byte

    // XX: this should take a signer as input
    // TagAndEncode() []byte
}

// TODO: create the right type of message here...
// CreateFromTLV(tlv codec.TLV) (Message, error)
func CreateFromTLV(tlv []codec.TLV) (Message, error) {
    var result Message

    return result, nil
}
