package context

type KeyTree struct {
    policyName string
    keys []Key
    Children []*KeyTree
}

func CreateKeyTree() (*KeyTree) {
    return nil, nil
}
