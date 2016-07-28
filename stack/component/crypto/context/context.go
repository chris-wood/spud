package context

import "fmt"

type CryptoContext struct {
    trustedKeys map[string]interface{}
}

type contextError struct {
    problem string
}

func (c contextError) Error() string {
    return fmt.Sprintf("%s", c.problem)
}

func NewCryptoContext() *CryptoContext {
    return &CryptoContext{
        trustedKeys: make(map[string]interface{}),
    }
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func (c *CryptoContext) GetTrustedKey(key string) (interface{}, bool) {
    val, ok := c.trustedKeys[key]
    return val, ok
}

func (c *CryptoContext) IsTrustedKey(key string) bool {
    _, ok := c.trustedKeys[key]
    return ok
}

func (c *CryptoContext) AddTrustedKey(key string, val interface{}) bool {
    _, ok := c.trustedKeys[key]
    if !ok {
        c.trustedKeys[key] = val
        return true
    }
    return false
}
