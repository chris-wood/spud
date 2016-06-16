package codec

import "encoding/binary"

type Encoder struct {
}

func (e Encoder) Encode(tlv TLV) []byte {
    nsType := make([]byte, 2)
    binary.BigEndian.PutUint16(nsType, tlv.Type())

    nsLength := make([]byte, 2)
    binary.BigEndian.PutUint16(nsLength, tlv.Length())

    tlTuple := append(nsType, nsLength...)

    return append(tlTuple, tlv.Value()...)
}
