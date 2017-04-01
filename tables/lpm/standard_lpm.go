package lpm

import "strings"

type prefixKey struct {
	value string
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
	if val, ok := t.dict[realKey]; ok {
		return val, nil
	}
	return nil, LPMError{"Not found"}
}

func (t *prefixTable) Drop(key string) {
	realKey := prefixKey{key}
	_, ok := t.dict[realKey]
	if ok {
		delete(t.dict, realKey)
	}
}

type StandardLPM struct {
	// XXX: we also need to store regexes here
	tables []*prefixTable
}

func (l *StandardLPM) extendTables(n int) {
	length := len(l.tables)
	if n > length {
		for i := 0; i < n-length; i++ {
			l.tables = append(l.tables, NewPrefixTable())
		}
	}
}

func (l *StandardLPM) Insert(keys []string, value interface{}) bool {
	l.extendTables(len(keys))

	for index, _ := range keys {
		prefix := strings.Join(keys[:index+1], "/")
		l.tables[index].Insert(prefix, value)
	}

	return true
}

func (l *StandardLPM) Lookup(keys []string) (interface{}, bool) {
	l.extendTables(len(keys))

	for index := len(keys) - 1; index >= 0; index-- {
		prefix := strings.Join(keys[:index+1], "/")
		if val, err := l.tables[index].Lookup(prefix); err == nil {
			return val, true
		}
	}
	return nil, false
}

func (l *StandardLPM) Drop(keys []string) {
	for index := len(keys) - 1; index >= 0; index-- {
		prefix := strings.Join(keys[:index+1], "/")
		if _, err := l.tables[index].Lookup(prefix); err == nil {
			l.tables[index].Drop(prefix)
		}
	}
}
