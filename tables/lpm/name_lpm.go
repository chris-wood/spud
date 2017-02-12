package lpm

import "github.com/chris-wood/spud/messages/name"

type NameLPM struct {
    // XXX
}

func (l *NameLPM) Insert(n name.Name, value interface{}) bool {
    // XXX: TODO
    return false
}

func (l *NameLPM) Lookup(n name.Name) (interface{}, bool) {
    // XXX: TODO
    return nil, false
}

func (l *NameLPM) Drop(n name.Name) {
    // XXX: TODO
}
