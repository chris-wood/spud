package messages

type Name struct {
    segments []NameSegment
}

func (n Name) Length() int {
    length := 0
    for _, ns := range(n.segments) {
        length += ns.Length()
    }
    return length + 4
}
