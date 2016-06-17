package main

import "fmt"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/codec"

func main() {
    ns1 := name.NameSegment{"foo"}
    ns2 := name.NameSegment{"bar"}

    name1 := name.Name{[]name.NameSegment{ns1, ns2}}
    name2, err := name.New("ccnx:/hello/spud")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(name2)

    fmt.Println(name1.Length());

    e := codec.Encoder{}
    fmt.Println(e.Encode(name1));
}
