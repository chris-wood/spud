package context

type KeyPath []Key

type KeyTree struct {
	PolicyName string
	Keys       []Key
	Children   []*KeyTree
}

func CreateKeyTree(policy string) *KeyTree {
	return &KeyTree{policy, nil, nil}
}

func (tree *KeyTree) AddKey(key Key) {
	tree.Keys = append(tree.Keys, key)
}

func (tree *KeyTree) AddChild(child *KeyTree) {
	tree.Children = append(tree.Children, child)
}

func (tree *KeyTree) GetKeyPaths() ([]KeyPath) {
    paths := make([]KeyPath, 0)

    // Do a DFS through the tree, accumulating the keys as you go
}
