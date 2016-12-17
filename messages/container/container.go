package container

import "fmt"

import "github.com/chris-wood/spud/codec"

type containerError struct {
    prob string
}

func (e containerError) Error() string {
    return fmt.Sprintf("%s", e.prob)
}

type Container interface {
    GetContainerType() uint16
    GetContainerValue() *codec.TLV
}
