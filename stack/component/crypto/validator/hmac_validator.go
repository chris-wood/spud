package validator

import (
	"hash"
	"crypto/hmac"
	"log"
	"crypto/sha256"
	"crypto/subtle"
	"github.com/chris-wood/spud/messages/validation"
	"github.com/chris-wood/spud/messages"
	ccnHash "github.com/chris-wood/spud/messages/hash"
	"github.com/chris-wood/spud/codec"
)

type HMACProcessor struct {
	secretKey []byte
	keyId ccnHash.Hash
}

func NewHMACProcessorWithKey(key []byte) (*HMACProcessor, error) {
	if len(key) != 16 && len(key) != 32 {
		return nil, processorError{"Invalid HMAC key length"}
	}

	// Create the KeyID
	hasher := sha256.New()
	hasher.Write(key)
	keyId := hasher.Sum(nil)
	keyIdHash := ccnHash.Create(ccnHash.HashTypeSHA256, keyId)

	return &HMACProcessor{secretKey: key, keyId: keyIdHash}, nil
}

func (p HMACProcessor) Sign(msg *messages.MessageWrapper) ([]byte, error) {
	mac := hmac.New(sha256.New, p.secretKey)
	digest := msg.HashProtectedRegion(sha256.New())
	mac.Write(digest)
	tag := mac.Sum(nil)
	return tag, nil
}

func (p HMACProcessor) Verify(request, response *messages.MessageWrapper) bool {
	validationPayload := response.GetValidationPayload()
	validationAlgorithm := response.GetValidationAlgorithm()

	if validationAlgorithm.GetValidationSuite() != codec.T_HMAC_SHA256 {
		log.Println("Invalid crypto type:", validationAlgorithm.GetValidationSuite())
		return false
	}

	// Compute the MAC
	mac := hmac.New(sha256.New, p.secretKey)
	digest := response.HashProtectedRegion(sha256.New())
	mac.Write(digest)
	tag := mac.Sum(nil)
	providedTag := validationPayload.Value()

	return subtle.ConstantTimeCompare(tag, providedTag) == 0
}

func (p HMACProcessor) ProcessorDetails() validation.ValidationAlgorithm {
	va := validation.NewValidationAlgorithmFromKeyId(codec.T_HMAC_SHA256, p.keyId, 0)
	return va
}

func (p HMACProcessor) Hasher() hash.Hash {
	return sha256.New()
}

func (p HMACProcessor) CanVerify(msg *messages.MessageWrapper) bool {
	validationAlgorithm := msg.GetValidationAlgorithm()

	if validationAlgorithm.GetValidationSuite() != codec.T_HMAC_SHA256 {
		return false
	}

	if !validationAlgorithm.GetKeyId().Equals(p.keyId) {
		return false
	}

	return true
}
