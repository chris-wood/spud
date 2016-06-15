package messages

type NameSegment struct {
    // TOOD: this really has a type and whatnot
    value string
}

func (ns NameSegment) Type() int {
    return 1
}

func (ns NameSegment) Length() int {
    return len(ns.value) + 4
}

func (ns NameSegment) Value() []byte {
    return []byte(ns.value)
}
