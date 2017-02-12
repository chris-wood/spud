package context

import "testing"

import "github.com/chris-wood/spud/messages/name"

func TestCreateSingleTree(t *testing.T) {
    rootName, _ := name.Parse("/foo")
    root := CreateKeyTree(*rootName)

    childNameA, _ := name.Parse("/foo/barA")
    childNameB, _ := name.Parse("/foo/barB")
    childA := CreateKeyTree(*childNameA)
    childB := CreateKeyTree(*childNameB)

    rootKey := SymmetricKey{[]byte{0x00, 0x01, 0x02, 0x03}}
    childKeyA := SymmetricKey{[]byte{0x00, 0x01, 0x02, 0x03}}
    childKeyB := SymmetricKey{[]byte{0x00, 0x01, 0x02, 0x03}}

    root.AddKey(rootKey)
    childA.AddKey(childKeyA)
    childB.AddKey(childKeyB)

    root.AddChild(childA)
    root.AddChild(childB)

    rootPaths := root.GetKeyPaths()
    if len(rootPaths) != 2 {
        t.Errorf("Invalid number of root keypaths")
    }
}
