package validator

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"github.com/chris-wood/spud/codec"
	"github.com/chris-wood/spud/messages"
	"github.com/chris-wood/spud/messages/validation"
	"github.com/chris-wood/spud/messages/validation/publickey"
	"hash"
	"log"
)

type RSAProcessor struct {
	privateKey *rsa.PrivateKey
	publicKey  rsa.PublicKey
}

func NewRSAProcessor(keySize int) (RSAProcessor, error) {
	var result RSAProcessor

	if keySize != 2048 && keySize != 4096 {
		return result, processorError{"Invalid key length provided: " + string(keySize)}
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return result, processorError{"Failed to generate a private key: " + err.Error()}
	}

	publicKey := privateKey.PublicKey

	return RSAProcessor{privateKey: privateKey, publicKey: publicKey}, nil
}

func NewRSAProcessorWithKey(key *rsa.PrivateKey) (RSAProcessor, error) {
	publicKey := key.PublicKey
	return RSAProcessor{privateKey: key, publicKey: publicKey}, nil
}

func (p RSAProcessor) Sign(msg *messages.MessageWrapper) ([]byte, error) {
	digest := msg.HashProtectedRegion(sha256.New())
	signature, err := rsa.SignPKCS1v15(rand.Reader, p.privateKey, crypto.SHA256, digest)
	return signature, err
}

func (p RSAProcessor) Verify(request, response *messages.MessageWrapper) bool {
	validationPayload := response.GetValidationPayload()
	validationAlgorithm := response.GetValidationAlgorithm()

	var key *rsa.PublicKey
	switch validationAlgorithm.GetValidationSuite() {
	case codec.T_RSA_SHA256:
		// XXX: the key might not be here...
		// we need a function that will, given a validation algorithm, resolve the key
		responseKey := validationAlgorithm.GetPublicKey()
		rawKey, err := x509.ParsePKIXPublicKey(responseKey.Value())
		if err != nil {
			log.Println("Error parsing public key")
			return false
		}
		key = rawKey.(*rsa.PublicKey)
	default:
		log.Println("Invalid crypto type:", validationAlgorithm.GetValidationSuite())
		return false
	}

	signature := validationPayload.Value()
	digest := response.HashProtectedRegion(sha256.New())
	err := rsa.VerifyPKCS1v15(key, crypto.SHA256, digest, signature)

	return err == nil
}

func (p RSAProcessor) ProcessorDetails() validation.ValidationAlgorithm {
	var result validation.ValidationAlgorithm
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&p.publicKey)
	if err != nil {
		return result
	}

	publicKey := publickey.Create(publicKeyBytes)
	va := validation.NewValidationAlgorithmFromPublickey(codec.T_RSA_SHA256, publicKey, 0)
	return va
}

func (p RSAProcessor) Hasher() hash.Hash {
	return sha256.New()
}

func (p RSAProcessor) CanVerify(msg *messages.MessageWrapper) bool {
	validationAlgorithm := msg.GetValidationAlgorithm()

	switch validationAlgorithm.GetValidationSuite() {
	case codec.T_RSA_SHA256:
		responseKey := validationAlgorithm.GetPublicKey()
		_, err := x509.ParsePKIXPublicKey(responseKey.Value())
		if err != nil {
			return false
		}
		return true
	default:
		log.Println("Invalid crypto type")
		return false
	}

	return false
}
