package main

import "fmt"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/name_segment"

func main() {

    ns1 := NameSegment{"foo"}
    ns2 := NameSegment{"bar"}

    name := Name{[]NameSegment{ns1, ns2}}

    fmt.Println(name.Length());
}
