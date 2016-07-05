package processor

import "github.com/chris-wood/spud/messages"
import "crypto/rand"
import "crypto/rsa"
import "crypto/sha256"
import "crypto"

type CryptoProcessor interface {
    Sign(msg messages.Message) ([]byte, error)
    Verify(msg messages.Message, signature []byte) bool
}

type RSAProcessor struct {
    privateKey *rsa.PrivateKey
    publicKey rsa.PublicKey
}

type processorError struct {
    problem string
}

func (p processorError) Error() string {
    return p.problem
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

func (p RSAProcessor) Sign(msg messages.Message) ([]byte, error) {
    digest := msg.HashSensitiveRegion(sha256.New())
    signature, err := rsa.SignPKCS1v15(rand.Reader, p.privateKey, crypto.SHA256, digest)
    return signature, err
}

func (p RSAProcessor) Verify(msg messages.Message, signature []byte) bool {
    digest := msg.HashSensitiveRegion(sha256.New())
    err := rsa.VerifyPKCS1v15(&p.publicKey, crypto.SHA256, digest, signature)
    return err != nil
}
