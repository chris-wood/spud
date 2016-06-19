package validation

import "fmt"

type ValidationAlgorithm struct {
    validationAlgorithmType uint16
    // TODO: validation dependent data goes here
}

type validationAlgorithmError struct {
    problem string
}

func (e validationAlgorithmError) Error() string {
    return fmt.Sprintf("%s", e.problem)
}

// Constructor functions

// TODO

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
