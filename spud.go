package main

import "fmt"
import "github.com/chris-wood/spud/messages/name"
// import "github.com/chris-wood/spud/messages/name_segment"
import "github.com/chris-wood/spud/codec"

func main() {
    // ns1 := name_segment.Parse("foo")
    // ns2 := name_segment.Parse("bar")

    // name1 := name.Name{[]name.NameSegment{ns1, ns2}}
    name2, err := name.Parse("ccnx:/hello/spud")
    if err != nil {
        fmt.Println(err)
    }
    // fmt.Println(name2)

    fmt.Println(name2.Length());

    e := codec.Encoder{}
    fmt.Println(e.Encode(name2));
}
