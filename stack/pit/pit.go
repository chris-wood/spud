package pit

import "github.com/chris-wood/spud/messages"

type PIT struct {
	table map[string]*messages.MessageWrapper
}

func NewPIT() *PIT {
	return &PIT{
		table: make(map[string]*messages.MessageWrapper),
	}
}

func (c *PIT) Insert(identity string, msg *messages.MessageWrapper) bool {
	_, ok := c.table[identity]
	if !ok {
		c.table[identity] = msg
		return true
	}
	return false
}

func (c *PIT) Lookup(identity string) (*messages.MessageWrapper, bool) {
	match, ok := c.table[identity]
	return match, ok
}

func (c *PIT) Remove(identity string) {
	_, ok := c.table[identity]
	if ok {
		delete(c.table, identity)
	}
}
