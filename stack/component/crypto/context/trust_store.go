package context

import (
	"github.com/chris-wood/spud/stack/component/crypto/context/schema"
	"fmt"
)

type TrustStore struct {
	trustedKeys map[string]interface{}
	trustSchema *schema.Schema
}

type contextError struct {
	problem string
}

func (c contextError) Error() string {
	return fmt.Sprintf("%s", c.problem)
}

func NewTrustStore() *TrustStore {
	return &TrustStore{
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

func (c *TrustStore) GetTrustedKey(key string) (interface{}, bool) {
	val, ok := c.trustedKeys[key]
	return val, ok
}

func (c *TrustStore) IsTrustedKey(key string) bool {
	_, ok := c.trustedKeys[key]
	return ok
}

func (c *TrustStore) AddTrustedKey(key string, val interface{}) bool {
	_, ok := c.trustedKeys[key]
	if !ok {
		c.trustedKeys[key] = val
		return true
	}
	return false
}