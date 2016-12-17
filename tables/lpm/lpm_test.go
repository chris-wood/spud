package lpm

import "testing"

import "github.com/chris-wood/spud/messages/name"

func createName(nameString string) name.Name {
    result, _ := name.Parse(nameString)
    return result
}

func TestInsert(t *testing.T) {
    var cases = []struct {
        inputName name.Name
        inserted bool
    }{
        { createName("/hello/world"), true},
        { createName("/hello/*"), true},
    }

    kvs := LPM{}
    value := 0

    for _, testCase := range cases {
        components := testCase.inputName.SegmentStrings()
        if testCase.inserted != kvs.Insert(components, value) {
            t.Errorf("Insert %s failed", components)
        }
    }
}

func TestLookup(t *testing.T) {
    // pass
}
