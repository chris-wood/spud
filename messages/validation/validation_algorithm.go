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

func createFromInnerTLV(validationType uint16, tlv codec.TLV) (ValidationAlgorithm, error) {
    var result ValidationAlgorithm
    var err error

    var publicKey publickey.PublicKey

    fmt.Println("parsing the validation alg innards")

    for _, child := range(tlv.Children()) {
        if child.Type() == codec.T_PUBLICKEY {
            publicKey, err = publickey.CreateFromTLV(child)
            if err != nil {
                return result, err
            }
        } else {
            fmt.Printf("Invalid TLV type %d\n", child.Type())
        }
    }

    return ValidationAlgorithm{validationAlgorithmType: validationType, publicKey: publicKey}, nil
}

func CreateFromTLV(tlv codec.TLV) (ValidationAlgorithm, error) {
    var result ValidationAlgorithm

    fmt.Println("parsing the validation alg")

    // There must be one child
    if len(tlv.Children()) != 1 {
        return result, nil
    }

    containerTlv := tlv.Children()[0]
    return createFromInnerTLV(containerTlv.Type(), containerTlv)
}

// ValidationAlgorithm functions

func (va ValidationAlgorithm) GetValidationSuite() uint16 {
    return va.validationAlgorithmType
}

func (va ValidationAlgorithm) GetPublicKey() publickey.PublicKey {
    return va.publicKey
}

// TLV interface functions

func (va ValidationAlgorithm) Type() uint16 {
    return codec.T_VALALG
}

func (va ValidationAlgorithm) TypeString() string {
    return "ValidationAlgorithm"
}

func (va ValidationAlgorithm) Length() uint16 {
    length := uint16(4) // 2+2 for the TL container

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

    // XXX: add the remaining VDD values here

    container := e.EncodeContainer(va.validationAlgorithmType, uint16(len(value)))
    container = append(container, value...)

    return container
}

func (va ValidationAlgorithm) Children() []codec.TLV  {
    children := []codec.TLV{va.publicKey}
    return children
}

func (va ValidationAlgorithm) String() string  {
    return va.TypeString()
}
