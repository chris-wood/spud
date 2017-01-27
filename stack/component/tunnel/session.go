package tunnel

import "crypto/aes"
import "crypto/cipher"

type Session struct {
    SessionID string
    counter int

    writeEncKey []byte
    writeMacKey []byte
    readEncKey []byte
    readMacKey []byte

    writeCipher cipher.AEAD
    readCipher cipher.AEAD
}

func NewSession(masterSecret []byte, sessionID string) *Session {
    WriteCipher, err := aes.NewCipher(masterSecret)
    if err != nil {
        panic(err.Error())
    }
    writeAEAD, err := cipher.NewGCM(WriteCipher)
    if err != nil {
        panic(err.Error())
    }

    ReadCipher, err := aes.NewCipher(masterSecret)
    if err != nil {
        panic(err.Error())
    }
    readAEAD, err := cipher.NewGCM(ReadCipher)
    if err != nil {
        panic(err.Error())
    }

    esic := Session {
        SessionID: sessionID,
        counter: 0,
        writeMacKey: masterSecret,
        writeEncKey: masterSecret,
        readMacKey: masterSecret,
        readEncKey: masterSecret,
        writeCipher: writeAEAD,
        readCipher: readAEAD,
    }

    return &esic
}

func (s Session) Encrypt(blob []byte) ([]byte, error) {
    nonce := make([]byte, 12)
    ciphertext := s.writeCipher.Seal(nil, nonce, blob, nil)
    return ciphertext, nil
}

func (s Session) Decrypt(blob []byte) ([]byte, error) {
    nonce := make([]byte, 12)
    return s.readCipher.Open(nil, nonce, blob, nil)
}
