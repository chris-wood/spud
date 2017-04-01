package validator

import "hash"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/validation"


type CryptoProcessor interface {
	CanVerify(msg *messages.MessageWrapper) bool
	Sign(msg *messages.MessageWrapper) ([]byte, error)
	Verify(request, response *messages.MessageWrapper) bool

	ProcessorDetails() validation.ValidationAlgorithm
	Hasher() hash.Hash
}


