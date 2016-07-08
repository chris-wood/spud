package validation

import "fmt"
import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages/link"
import "github.com/chris-wood/spud/messages/validation/publickey"

type ValidationAlgorithm struct {
    validationAlgorithmType uint16

    // Validation dependent data -- empty until otherwise instantiated
    publicKey publickey.PublicKey
    signatureTime uint64

    // XXX: write TLV wrappers for these fields
    keyId []byte
    certificate []byte
    keyName link.Link
}

type validationAlgorithmError struct {
    problem string
}

func (e validationAlgorithmError) Error() string {
    return fmt.Sprintf("%s", e.problem)
}

// Constructor functions

// func NewValidationAlgorithm(vaType uint16, keyId, publicKey, certificate []byte, keyName link.Link, signatureTime uint64) ValidationAlgorithm {
//     return ValidationAlgorithm{validationAlgorithmType: vaType, keyId: keyId, publicKey: publicKey, certificate: certificate, keyName: keyName, signatureTime: signatureTime}
// }

func NewValidationAlgorithmFromPublickey(vaType uint16, publicKey publickey.PublicKey, signatureTime uint64) ValidationAlgorithm {
    return ValidationAlgorithm{validationAlgorithmType: vaType, publicKey: publicKey, signatureTime: signatureTime}
}

// func NewValidationAlgorithmFromKeyId(vaType uint16, keyId []byte, signatureTime uint64) ValidationAlgorithm {
//     return ValidationAlgorithm{validationAlgorithmType: vaType, keyId: keyId, signatureTime: signatureTime}
// }
//
// func NewValidationAlgorithmFromLink(vaType uint16, keyName link.Link, signatureTime uint64) ValidationAlgorithm {
//     return ValidationAlgorithm{validationAlgorithmType: vaType, keyName: keyName, signatureTime: signatureTime}
// }
//
// func NewValidationAlgorithmFromCertificate(vaType uint16, certificate []byte, signatureTime uint64) ValidationAlgorithm {
//     return ValidationAlgorithm{validationAlgorithmType: vaType, certificate: certificate, signatureTime: signatureTime}
// }

func CreateFromTLV(tlv codec.TLV) (ValidationAlgorithm, error) {
    var result ValidationAlgorithm
    var err error

    var publicKey publickey.PublicKey

    for _, child := range(tlv.Children()) {
        if child.Type() == codec.T_PUBLICKEY {
            publicKey, err = publickey.CreateFromTLV(child)
            if err != nil {
                return result, err
            }
        }
    }

    return ValidationAlgorithm{validationAlgorithmType: tlv.Type(), publicKey: publicKey}, nil
}

// ValidationAlgorithm functions

func (va ValidationAlgorithm) GetPublicKey() publickey.PublicKey {
    return va.publicKey
}

// TLV interface functions

func (va ValidationAlgorithm) Type() uint16 {
    return va.validationAlgorithmType
}

func (va ValidationAlgorithm) TypeString() string {
    return "ValidationAlgorithm"
}

func (va ValidationAlgorithm) Length() uint16 {
    length := uint16(0)

    if va.publicKey.Length() > 0  {
        length += va.publicKey.Length() + 4
    }

    // XXX: add the remaining values here

    return length
}

func (va ValidationAlgorithm) Value() []byte {
    e := codec.Encoder{}
    value := make([]byte, 0)

    if va.publicKey.Length() > 0  {
        value = append(value, e.EncodeTLV(va.publicKey)...)
    }

    // XXX: add the remaining values here

    return value
}

func (va ValidationAlgorithm) Children() []codec.TLV  {
    children := []codec.TLV{}
    return children
}

func (va ValidationAlgorithm) String() string  {
    return va.TypeString()
}
