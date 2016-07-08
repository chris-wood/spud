package codec

import "encoding/binary"

type Encoder struct {
}

func (e Encoder) EncodeContainer(tlvType, tlvLength uint16) []byte {
    nsType := make([]byte, 2)
    binary.BigEndian.PutUint16(nsType, tlvType)

    nsLength := make([]byte, 2)
    binary.BigEndian.PutUint16(nsLength, tlvLength)
    tlTuple := append(nsType, nsLength...)

    return tlTuple
}

func (e Encoder) EncodeTLV(tlv TLV) []byte {
    nsType := make([]byte, 2)
    binary.BigEndian.PutUint16(nsType, tlv.Type())

    nsLength := make([]byte, 2)
    binary.BigEndian.PutUint16(nsLength, tlv.Length())

    tlTuple := append(nsType, nsLength...)
    return append(tlTuple, tlv.Value()...)
}

func (e Encoder) Encode(tlvList []TLV) []byte {
    encodedBytes := make([]byte, 0)
    for _, tlv := range(tlvList) {
        encodedBytes = append(encodedBytes, e.EncodeTLV(tlv)...)
    }
    return encodedBytes
}
