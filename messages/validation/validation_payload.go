package validation

import "fmt"
import "github.com/chris-wood/spud/codec"

type ValidationPayload struct {
    bytes []byte
}

type validationPayloadError struct {
    problem string
}

func (e validationPayloadError) Error() string {
    return fmt.Sprintf("%s", e.problem)
}

// Constructor functions

func NewValidationPayload(bytes []byte) ValidationPayload {
    return ValidationPayload{}
}

// TLV interface functions

func (va ValidationPayload) Type() uint16 {
    return codec.T_VALPAYLOAD
}

func (va ValidationPayload) TypeString() string {
    return "ValidationPayload"
}

func (va ValidationPayload) Length() uint16 {
    return uint16(len(va.bytes))
}

func (va ValidationPayload) Value() []byte {
    return va.bytes
}

// String functions

func (va ValidationPayload) String() string {
    return ""
}
