package lpm

import "testing"
import "github.com/chris-wood/spud/messages/name"

func TestInsert(t *testing.T) {
    names := []name.Name{
        name.Parse("ccnx:/hello/world"),
        name.Parse("ccnx:/hello/*"),
    }

    kvs := LPM{}

    for i, n := range(names) {
        components := n.SegmentStrings()
        if kvs.Insert(components) != true {
            t.Errorf("Insert %s failed", components)
        }

        // XXX: assert that the number of entries actually increased
    }
}

func TestLookup(t *testing.T) {
    // pass
}
