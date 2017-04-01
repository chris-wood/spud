package context

import "github.com/chris-wood/spud/tables/lpm"
import "github.com/chris-wood/spud/messages/name"

import "fmt"

type KeyPath []Key

type KeyStore struct {
	keyTree lpm.LPM
}

type keyStoreError struct {
	problem string
}

func (c keyStoreError) Error() string {
	return fmt.Sprintf("%s", c.problem)
}

func NewKeyStore() *KeyStore {
	return &KeyStore{
		keyTree: &lpm.StandardLPM{},
	}
}

func (ke *KeyStore) AddKey(nameSchema name.Name, theKey Key) {
	segments := nameSchema.SegmentStrings()
	if keyTreeBlob, ok := ke.keyTree.Lookup(segments); ok {
		tree := keyTreeBlob.(*KeyTree)
		tree.AddKey(theKey)
	} else {
		tree := CreateKeyTree(nameSchema)
		tree.AddKey(theKey)
		ke.keyTree.Insert(segments, tree)
	}
}

func (ke *KeyStore) GetKeys(nameSchema name.Name) ([]KeyPath, error) {
	segments := nameSchema.SegmentStrings()
	if keyTreeBlob, ok := ke.keyTree.Lookup(segments); ok {
		keyTree := keyTreeBlob.(*KeyTree)
		paths := keyTree.GetKeyPaths()
		return paths, nil
	} else {
		return []KeyPath{}, keyStoreError{"Key for name not found"}
	}
}
