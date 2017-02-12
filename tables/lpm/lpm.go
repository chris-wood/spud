package lpm

type LPMError struct {
    message string
}

func (e LPMError) Error() string {
    return e.message
}

type LPM interface {
    Insert(keys []string, value interface{}) bool
    Lookup(keys []string) (interface{}, bool)
    Drop(keys []string)
}
