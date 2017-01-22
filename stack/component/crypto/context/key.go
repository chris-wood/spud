package context

// Purpose
const KeyPurposeEncrypt uint8 = 0
const KeyPurposeSign uint8 = 1

// KeyType
// XXX

type Key struct {
    Purpose uint8
    KeyType uint16
    Value interface{}
}
