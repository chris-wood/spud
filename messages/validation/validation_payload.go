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

func (vp ValidationPayload) Type() uint16 {
    return codec.T_VALPAYLOAD
}

func (vp ValidationPayload) TypeString() string {
    return "ValidationPayload"
}

func (vp ValidationPayload) Length() uint16 {
    return uint16(len(vp.bytes))
}

func (vp ValidationPayload) Value() []byte {
    return vp.bytes
}

func (vp ValidationPayload) Children() []codec.TLV {
    return nil
}

// String functions

func (vp ValidationPayload) String() string {
    return ""
}
