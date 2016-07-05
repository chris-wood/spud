package processor

import "github.com/chris-wood/spud/messages"
import "crypto"
import "crypto/rand"
import "crypto/rsa"
import "crypto/sha256"

type CryptoProcessor interface {
    Sign(msg messages.Message) []byte
    Verify(msg messages.Message) bool
}

type RSAProcessor struct {
    privateKey *rsa.PrivateKey
    publicKey *rsa.PublicKey
}

type processorError struct {
    problem string
}

func (p processorError) Error() string {
    return problem
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

    return RSAProcessor{privateKey: privateKey, publicKey: publicKey}
}

func (p RSAProcessor) Sign(msg messages.Message) []byte {
    digest := msg.HashSensitiveRegion()

}

func (p RSAProcessor) Verify(msg messages.Message) bool {

}
