package context

import "testing"

import "github.com/chris-wood/spud/messages/name"

func TestKeyStore(t *testing.T) {
    store := NewKeyStore()

    rootName, _ := name.Parse("/foo")
    childNameA, _ := name.Parse("/foo/barA")
    childNameB, _ := name.Parse("/foo/barB")

    rootKey := SymmetricKey{[]byte{0x00, 0x01, 0x02, 0x03}}
    childKeyA := SymmetricKey{[]byte{0x00, 0x01, 0x02, 0x03}}
    childKeyB := SymmetricKey{[]byte{0x00, 0x01, 0x02, 0x03}}

    store.AddKey(*rootName, rootKey)
    store.AddKey(*childNameA, childKeyA)
    store.AddKey(*childNameB, childKeyB)

    rootPaths, err := store.GetKeys(*childNameA)
    if err != nil {
        t.Errorf("An error occurred when fetching the key: %s", err.Error())
    }
    if len(rootPaths) != 1 {
        t.Errorf("Invalid key path for childNameA. Got %d, expected 1", len(rootPaths))
    }

    // root := CreateKeyTree(*rootName)
    //
    // childNameA, _ := name.Parse("/foo/barA")
    // childNameB, _ := name.Parse("/foo/barB")
    // childA := CreateKeyTree(*childNameA)
    // childB := CreateKeyTree(*childNameB)
    //
    // rootKey := SymmetricKey{[]byte{0x00, 0x01, 0x02, 0x03}}
    // childKeyA := SymmetricKey{[]byte{0x00, 0x01, 0x02, 0x03}}
    // childKeyB := SymmetricKey{[]byte{0x00, 0x01, 0x02, 0x03}}
    //
    // root.AddKey(rootKey)
    // childA.AddKey(childKeyA)
    // childB.AddKey(childKeyB)
    //
    // root.AddChild(childA)
    // root.AddChild(childB)
    //
    // rootPaths := root.GetKeyPaths()
    // if len(rootPaths) != 2 {
    //     t.Errorf("Invalid number of root keypaths")
    // }
}
