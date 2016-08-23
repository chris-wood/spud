package pit

import "github.com/chris-wood/spud/messages"

type PIT struct {
    table map[string]messages.Message
}

func NewPIT() *PIT {
    return &PIT{
        table: make(map[string][]byte),
    }
}

func (c *PIT) Insert(identity string, msg messages.Message) bool {
    _, ok := c.table[identity]
    if !ok {
        c.table[identity] = msg
        return true
    }
    return false
}

func (c *PIT) Lookup(identity string) (messages.Message, bool) {
    match, ok := c.table[identity]
    return match, ok
}

func (c *PIT) Remove(identity string) {
    match, ok := c.table[identity]
    if ok {
        delete(c, identity)
    }
}
