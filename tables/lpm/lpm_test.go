package lpm

import "testing"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/tables/lpm"

func createName(nameString string) name.Name {
    result, _ := name.Parse(nameString)
    return result
}

func TestInsert(t *testing.T) {
    names := []name.Name{
        createName("ccnx:/hello/world"),
        createName("ccnx:/hello/*"),
    }

    kvs := lpm.LPM{}
    value := 0

    for _, n := range(names) {
        components := n.SegmentStrings()
        if kvs.Insert(components, value) != true {
            t.Errorf("Insert %s failed", components)
        }

        // XXX: assert that the number of entries actually increased
    }
}

func TestLookup(t *testing.T) {
    // pass
}
