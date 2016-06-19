package messages

import "github.com/chris-wood/spud/codec"

type Message interface {
    CreateFromTLV(tlv codec.TLV) *Message
    ComputeMessageHash() []byte
}
