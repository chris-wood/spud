package main

import "os"
import "fmt"
import "bufio"
import "crypto/rand"
import "crypto/rsa"
import "crypto/x509"
import "encoding/pem"

type keyError struct {
	problem string
}

func (k keyError) Error() string {
	return fmt.Sprintf("%s", k.problem)
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	// case *ecdsa.PrivateKey:
	//     b, err := x509.MarshalECPrivateKey(k)
	//     if err != nil {
	//         fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
	//         os.Exit(2)
	//     }
	//     return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func generateKeyPair(keySize int) (*rsa.PrivateKey, error) {
	var result *rsa.PrivateKey

	if keySize != 2048 && keySize != 4096 {
		return result, keyError{"Invalid key length provided: " + string(keySize)}
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return result, keyError{"Failed to generate a private key: " + err.Error()}
	}

	return privateKey, nil
}

func main() {
	sk, err := generateKeyPair(2048)
	if err != nil {
		os.Exit(-1)
	}

	keyOut := bufio.NewWriter(os.Stdout)
	defer keyOut.Flush()
	pem.Encode(keyOut, pemBlockForKey(sk))
}
