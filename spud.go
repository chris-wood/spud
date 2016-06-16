package main

import "fmt"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/codec"

func main() {
    ns1 := messages.NameSegment{"foo"}
    ns2 := messages.NameSegment{"bar"}

    name := messages.Name{[]messages.NameSegment{ns1, ns2}}

    fmt.Println(name.Length());

    e := codec.Encoder{}
    fmt.Println(e.Encode(name));
}
