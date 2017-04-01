package context

// Key Type
const KeyTypeEC uint8 = 0
const KeyTypeRSA uint8 = 1
const KeyTypeOpaque uint8 = 2

// Curves
const KeyCurveP256 uint8 = 0

// ...

// Encryption algorithms
const KeyEncryptionAlgorithm_AES_128 uint16 = 0
const KeyEncryptionAlgorithm_AES_256 uint16 = 1
const KeyEncryptionAlgorithm_RSA_2048_OAEP uint16 = 2
const KeyEncryptionAlgorithm_RSA_4096_OAEP uint16 = 3

type Key interface {
	Purpose() uint8
	KeyType() uint8
	Algorithm() uint16

	Value() []byte
}

type SymmetricKey struct {
	bytes []byte
}

func (k SymmetricKey) Purpose() uint8 {
	return uint8(0)
}

func (k SymmetricKey) KeyType() uint8 {
	return KeyTypeOpaque
}

func (k SymmetricKey) Algorithm() uint16 {
	return uint16(0)
}

func (k SymmetricKey) Value() []byte {
	return k.bytes
}
