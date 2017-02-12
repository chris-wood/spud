package context

import "github.com/chris-wood/spud/messages/name"

type KeyTree struct {
	PolicyName name.Name
	Keys       []Key
	Children   []*KeyTree
}

func CreateKeyTree(policyName name.Name) *KeyTree {
	return &KeyTree{policyName, nil, nil}
}

func (tree *KeyTree) AddKey(key Key) {
	tree.Keys = append(tree.Keys, key)
}

func (tree *KeyTree) AddChild(child *KeyTree) {
	tree.Children = append(tree.Children, child)
}

func (tree *KeyTree) getKeyPathsWithAccumulator(acc KeyPath) (KeyPath) {
    // paths := make([][]KeyPath, 0)
    // XXX: TODO
    return acc
}

func (tree *KeyTree) GetKeyPaths() ([]KeyPath) {
    paths := make([]KeyPath, 0)

    // rootKeys := make([]KeyPath, 0)
    // rootKeys = append(rootKeys, tree.Keys)
    rootKeys := tree.Keys[:]

    for _, child := range(tree.Children) {
        subpath := child.getKeyPathsWithAccumulator(rootKeys[:])
        paths = append(paths, subpath)
    }

    return paths
}
