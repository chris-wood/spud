package validator

// import "github.com/chris-wood/spud/codec"
// import "github.com/chris-wood/spud/messages"
// import "github.com/chris-wood/spud/messages/validation"
// import "github.com/chris-wood/spud/messages/validation/publickey"

// import "crypto/rand"
// import "crypto/rsa"
// import "crypto/sha256"
// import "crypto/x509"
// import "crypto"
// import "hash"
// import "log"
import "fmt"

type Encryptor interface {
    Encrypt(identifier string, payload []byte) ([]byte, error)
    Decrypt(identifier string, payload []byte) ([]byte, error)
}

type RSAEncryptor struct {
    // privateKey *rsa.PrivateKey
    // publicKey rsa.PublicKey
}

type processorError struct {
    problem string
}

func (p processorError) Error() string {
    return fmt.Sprintf("%s", p.problem)
}

// func NewRSAEncryptor() (RSAEncryptor, error) {
    // var result RSAProcessor
    //
    // if keySize != 2048 && keySize != 4096 {
    //     return result, processorError{"Invalid key length provided: " + string(keySize)}
    // }
    //
    // privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
    // if err != nil {
    //     return result, processorError{"Failed to generate a private key: " + err.Error()}
    // }
    //
    // publicKey := privateKey.PublicKey
    //
    // return RSAProcessor{privateKey: privateKey, publicKey: publicKey}, nil
// }

// func (p RSAEncryptor) Encrypt(identifier string, payload []byte) ([]byte, error) {
//     // XXX
// }
//
// func (p RSAEncryptor) Decrypt(identifier string, payload []byte) bool {
//     // XXX
// }
