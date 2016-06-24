package crypto

import "github.com/chris-wood/spud/messages"

type CryptoProcessor interface {
    Sign(msg *messages.Message) []byte
    Verify(msg *messages.Message) bool
}

type XXXProcessor struct {
    // XX key store
}

func (c Crypto) ProcessEgressMessages() {
    for ;; {
        // XX

        // 1. if sign, then compute hash, sign, and append validation alg/locator
        // 2. if mac, then compute hash, mac, and append stuff
    }
}

func (c Crypto) ProcessIngressMessages() {
    for ;; {
        // XX
    }
}
