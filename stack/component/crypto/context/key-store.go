package context

import "github.com/chris-wood/spud/tables/lpm"
import "fmt"

type KeyStore struct {
	keyTree *lpm.LPM
}

type trustStoreError struct {
	problem string
}

func (c trustStoreError) Error() string {
	return fmt.Sprintf("%s", c.problem)
}

func NewKeyStore() *KeyStore {
	return &KeyStore{
		keyTree: new(lpm.LPM),
	}
}

func (ke *KeyStore) AddKey(nameSchema string, theKey Key) {
	// XXX: using the schema, identify the key tree
	// if one exists, insert the key into that tree
	// else, create a new key tree and add it to the root
}
