package validation

import "fmt"
import "github.com/chris-wood/spud/messages/link"

type ValidationAlgorithm struct {
    validationAlgorithmType uint16

    // Validation dependent data -- empty until otherwise instantiated
    keyId []byte
    publicKey []byte
    certificate []byte
    keyName link.Link
    signatureTime uint64
}

type validationAlgorithmError struct {
    problem string
}

func (e validationAlgorithmError) Error() string {
    return fmt.Sprintf("%s", e.problem)
}

// Constructor functions

func NewValidationAlgorithm(vaType uint16, keyId, publicKey, certificate []byte, keyName link.Link, signatureTime uint64) ValidationAlgorithm {
    return ValidationAlgorithm{validationAlgorithmType: vaType, keyId: keyId, publicKey: publicKey, certificate: certificate, keyName: keyName, signatureTime: signatureTime}
}

func NewValidationAlgorithmFromPublickey(vaType uint16, publicKey []byte, signatureTime uint64) ValidationAlgorithm {
    return ValidationAlgorithm{validationAlgorithmType: vaType, publicKey: publicKey, signatureTime: signatureTime}
}

func NewValidationAlgorithmFromKeyId(vaType uint16, keyId []byte, signatureTime uint64) ValidationAlgorithm {
    return ValidationAlgorithm{validationAlgorithmType: vaType, keyId: keyId, signatureTime: signatureTime}
}

func NewValidationAlgorithmFromLink(vaType uint16, keyName link.Link, signatureTime uint64) ValidationAlgorithm {
    return ValidationAlgorithm{validationAlgorithmType: vaType, keyName: keyName, signatureTime: signatureTime}
}

func NewValidationAlgorithmFromCertificate(vaType uint16, certificate []byte, signatureTime uint64) ValidationAlgorithm {
    return ValidationAlgorithm{validationAlgorithmType: vaType, certificate: certificate, signatureTime: signatureTime}
}

// TLV interface functions

func (va ValidationAlgorithm) Type() uint16 {
    return va.validationAlgorithmType
}

func (va ValidationAlgorithm) TypeString() string {
    return "ValidationAlgorithm"
}

func (va ValidationAlgorithm) Length() uint16 {
    return 0
}

func (va ValidationAlgorithm) Value() []byte {
    return make([]byte, 0)
}

// String functions

func (va ValidationAlgorithm) String() string {
    return ""
}
