package main

import "fmt"

type TLV interface {
    Type() int
    Length() int
    Value() []byte
}

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

func main() {
//    var x map[string]string
    
    ns1 := NameSegment{"foo"}
    ns2 := NameSegment{"bar"}

    name := Name{[]NameSegment{ns1, ns2}}

    fmt.Println(name.Length());
}
