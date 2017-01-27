package context

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
