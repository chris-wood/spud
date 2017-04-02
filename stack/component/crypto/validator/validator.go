package validator

import (
	"github.com/chris-wood/spud/messages/validation"
	"github.com/chris-wood/spud/messages"
	"fmt"
	"hash"
)

type processorError struct {
	problem string
}

func (p processorError) Error() string {
	return fmt.Sprintf("%s", p.problem)
}

type CryptoProcessor interface {
	CanVerify(msg *messages.MessageWrapper) bool
	Sign(msg *messages.MessageWrapper) ([]byte, error)
	Verify(request, response *messages.MessageWrapper) bool

	ProcessorDetails() validation.ValidationAlgorithm
	Hasher() hash.Hash
}
