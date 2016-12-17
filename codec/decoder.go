package codec

import "fmt"
import "encoding/binary"

type Decoder struct {
}

type decoderError struct {
    problem string
}

func (e decoderError) Error() string {
    return fmt.Sprintf("%s", e.problem)
}

func hasInnerTLV(tlvType, tlvLength uint16, bytes []byte) bool {
    if tlvLength < 4 {
        return false
    }

    subLength := readWord(bytes[2:])

    if subLength < tlvLength {
        return true
    } else {
        return false
    }
}

func readWord(bytes []byte) uint16 {
    return binary.BigEndian.Uint16(bytes)
}

func (d Decoder) decodeTLV(tlvType, tlvLength uint16, bytes []byte) TLV {
    if hasInnerTLV(tlvType, tlvLength, bytes) {
        children := make([]TLV, 0)
        for offset := uint16(0); offset < tlvLength && offset < uint16(len(bytes)); {
            innerType := readWord(bytes[offset:])
            offset += 2
            innerLength := readWord(bytes[offset:])
            offset += 2

            if offset + innerLength > tlvLength { // Failure.
                return NewLeafTLV(tlvType, bytes[0:tlvLength])
            } else {
                child := d.decodeTLV(innerType, innerLength, bytes[offset:(offset + innerLength)])
                children = append(children, child)
            }

            offset += innerLength
        }

        tlv := NewNestedTLV(tlvType, children)
        return tlv
    } else {
        return NewLeafTLV(tlvType, bytes[0:tlvLength])
    }
}

func (d Decoder) Decode(bytes []byte) []TLV {
    tlvs := make([]TLV, 0)

    for index := 0; index < len(bytes); {
        tlvType := readWord(bytes[index:])
        tlvLength := readWord(bytes[index + 2:])

        tlv := d.decodeTLV(tlvType, tlvLength, bytes[index + 4:])
        tlvs = append(tlvs, tlv)

        index += 4 + int(tlvLength)
    }

    return tlvs
}
