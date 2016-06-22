package link

import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/hash"
import "fmt"

// import "encoding/json"

type Link struct {
    linkName *name.Name
    keyId *hash.Hash
    contentId *hash.Hash
}

type linkError struct {
    prob string
}

func (e linkError) Error() string {
    return fmt.Sprintf("%s", e.prob)
}
