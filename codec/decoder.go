package codec

// import "encoding/binary"
import "fmt"

type Decoder struct {
}

// func (d Decoder) DecodeValidationAlgorithmTLV(bytes []byte) TLVInterface {
//     // tlvType := readWord(bytes[0:])
//     // tlvLength := readWord(bytes[2:])
//
//     switch tlvType {
//     case T_CRC32C:
//     case T_HMAC_SHA256:
//     case T_RSA_SHA256:
//     case T_KEYID:
//     case T_PUBLICKEY:
//     case T_SIGTIME:
//     }
//     return nil
// }
//
// func (d Decoder) DecodeValidationPayloadTLV(ytes []byte) TLVInterface {
//     return nil
// }
//
// func (d Decoder) DecodePayloadTLV(bytes []byte) TLVInterface {
//     // tlvType := readWord(bytes[0:])
//     // tlvLength := readWord(bytes[2:])
//
//     switch tlvType {
//     case T_PAYLOADTYPE_DATA:
//     case T_PAYLOADTYPE_KEY:
//     case T_PAYLOADTYPE_LINK:
//     case T_PAYLOADTYPE_MANIFEST:
//     }
//     return nil
// }
//
// func (d Decoder) DecodeNameTLV(bytes []byte) TLVInterface {
//     // tlvType := readWord(bytes[0:])
//     // tlvLength := readWord(bytes[2:])
//
//     switch tlvType {
//     case T_NAMESEG_NAME:
//     case T_NAMESEG_IPID:
//     case T_NAMESEG_CHUNK:
//     case T_NAMESEG_VERSION:
//     case T_NAMESEG_APP0:
//     case T_NAMESEG_APP1:
//     case T_NAMESEG_APP2:
//     case T_NAMESEG_APP3:
//     case T_NAMESEG_APP4:
//     }
//     return nil
// }
//
// func (d Decoder) DecodeMessageTLV(bytes []byte) TLVInterface {
//     tlvType := readWord(bytes[0:])
//     tlvLength := readWord(bytes[2:])
//
//     switch tlvType {
//     case T_NAME:
//
//     case T_PAYLOAD:
//     case T_KEYID_REST:
//     case T_HASH_REST:
//     case T_PAYLDTYPE:
//     case T_EXPIRY:
//     case T_HASHGROUP:
//     case T_BLOCKHASHGROUP:
//     }
//     return nil
// }
//
// func (d Decoder) DecodeTopLevelTLV(bytes []byte) TLVInterface {
//     // tlvType := readWord(bytes[0:])
//     // tlvLength := readWord(bytes[2:])
//
//     switch tlvType {
//     case T_INTEREST:
//     case T_OBJECT:
//     case T_VALALG:
//     case T_VALSIG:
//     case T_MANIFEST:
//     }
//     return nil
// }
//
// func (d Decoder) DecodeHopByHopTLV(bytes []byte) TLVInterface {
//     // tlvType := readWord(bytes[0:])
//     // tlvLength := readWord(bytes[2:])
//
//     switch tlvType {
//     case T_INT_LIFE:
//     case T_CACHE_TIME:
//     case T_MSG_HASH:
//     }
//     return nil
// }

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

func (d Decoder) decodeTLV(tlvType, tlvLength uint16, bytes []byte) TLVInterface {
    if hasInnerTLV(tlvType, tlvLength, bytes) {
        children := make([]TLVInterface, 0)
        for offset := uint16(4); offset < tlvLength; {
            innerType := readWord(bytes[offset:])
            offset += 2
            innerLength := readWord(bytes[offset:])
            offset += 2

            if offset + innerLength > tlvLength { // Failure.
                return NewLeafTLV(tlvType, bytes[0:tlvLength])
            } else {
                child := d.decodeTLV(innerType, innerLength, bytes[offset:])
                children = append(children, child)
            }

            offset += innerLength
        }

        return NewNestedTLV(children)
    } else {
        return NewLeafTLV(tlvType, bytes[0:tlvLength])
    }
}

func readWord(bytes []byte) uint16 {
    // binary.BigEndian.PutUint16(nsType, tlv.Type())
    return uint16(bytes[0] << 8 | bytes[1])
}

func (d Decoder) Decode(bytes []byte) []TLVInterface {
    tlvs := make([]TLVInterface, 0)

    for index := 0; index < len(bytes); {
        tlvType := readWord(bytes[index:])
        tlvLength := readWord(bytes[index + 2:])

        tlv := d.decodeTLV(tlvType, tlvLength, bytes)
        tlvs = append(tlvs, tlv)

        index += 4 + int(tlvLength)
    }

    return tlvs
}
