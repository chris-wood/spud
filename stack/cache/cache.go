package cache

type Cache struct {
    table map[string][]byte
}

func NewCache() *Cache {
    return &Cache{
        table: make(map[string][]byte),
    }
}

func (c *Cache) Insert(identity string, wireFormat []byte) bool {
    _, ok := c.table[identity]
    if !ok {
        c.table[identity] = wireFormat
        return true
    } else { // aggregate
        return false
    }
}

func (c *Cache) Lookup(identity string) ([]byte, bool) {
    match, ok := c.table[identity]
    return match, ok
}
