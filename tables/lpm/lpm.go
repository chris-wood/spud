package lpm

import "strings"

type prefixKey struct {
    value string
}

type lpmError struct {
    message string
}

func (e lpmError) Error() string {
    return e.message
}

func (pk prefixKey) Compare(other prefixKey) bool {
    return false
}

type prefixTable struct {
    dict map[prefixKey]interface{}
}

func NewPrefixTable() *prefixTable {
    return &prefixTable{make(map[prefixKey]interface{})}
}

func (t *prefixTable) Insert(key string, value interface{}) {
    realKey := prefixKey{key}
    t.dict[realKey] = value
}

func (t *prefixTable) Lookup(key string) (interface{}, error) {
    realKey := prefixKey{key}
    val, ok := t.dict[realKey]
    if !ok {
        return val, lpmError{"Not found"}
    }
    return val, nil
}

type LPM struct {
    // XXX: we also need to store regexes here
    tables []*prefixTable
}

func (l *LPM) extendTables(n int) {
    length := len(l.tables)
    if n > length {
        for i := 0; i < n - length ; i++ {
            l.tables = append(l.tables, NewPrefixTable())
        }
    }
}

func (l *LPM) Insert(keys []string, value interface{}) bool {
    l.extendTables(len(keys))

    for index, _ := range(keys) {
        prefix := strings.Join(keys[:index + 1], "")
        l.tables[index].Insert(prefix, value)
    }

    return true
}

func (l *LPM) Lookup(keys []string) (interface{}, bool) {
    l.extendTables(len(keys))

    // XXX: this is wrong -- we need to lookup from the
    for index, _ := range(keys) {
        prefix := strings.Join(keys[:index + 1], "")
        if val, err := l.tables[index].Lookup(prefix); err == nil {
            return val, true
        }
    }
    return nil, false
}

func (l LPM) Drop(key string) {
}
