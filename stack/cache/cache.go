package cache

import "github.com/chris-wood/spud/messages"

type Cache struct {
    table map[string]messages.Message
}

func NewCache() *CodecComponent {
    return &Cache{
        table: make(map[string]messages.Message)
    }
}

func (c *Cache) Insert(msg messages.Message) bool {
    identity := msg.Identifier()
    _, ok := c.table[identity]
    if !ok {
        // XXX: apply eviction strategy here...
        c.table[identity] = msg
        return true
    }
    return false
}

func (c *Cache) Lookup(msg messages.Message) (messages.Message, bool) {
    identity := msg.Identifier()
    match, ok := c.table[identity]
    return match, ok
}
