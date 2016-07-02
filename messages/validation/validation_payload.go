package validation

import "fmt"

type ValidationPayload struct {
    // TODO
}

type validationPayloadError struct {
    problem string
}

func (e validationPayloadError) Error() string {
    return fmt.Sprintf("%s", e.problem)
}

// Constructor functions

// TODO

// TLV interface functions

func (va ValidationPayload) Type() uint16 {
    return 0
}

func (va ValidationPayload) TypeString() string {
    return "ValidationPayload"
}

func (va ValidationPayload) Length() uint16 {
    return 0
}

func (va ValidationPayload) Value() []byte {
    return make([]byte, 0)
}

// String functions

func (va ValidationPayload) String() string {
    return ""
}
