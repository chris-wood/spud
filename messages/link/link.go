package link

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

// Constructors

func Create(linkName *name.Name, keyId, contentId *hash.Hash) *Link {
    return &Link{linkName, keyId, contentId}
}

// API

func (l Link) Name() *name.Name {
    return l.linkName
}

func (l Link) KeyID() *hash.Hash {
    return l.keyId
}

func (l Link) ContentID()  *hash.Hash {
    return l.contentId
}

// // TLV functions
//
// func (l Link) Type() uint16 {
//     return uint16(codec.T_LINK)
// }
//
// func (l Link) TypeString() string {
//     return "Link"
// }
//
// func (l Link) Length() uint16 {
//     return len(l.Value())
// }
//
// func (l Link) Value() []byte  {
//     bytes := l.linkName.Value()
//     bytes = append(bytes, l.keyId.Value())
//     bytes = append(bytes, l.contentId.Value())
//     return bytes
// }
//
// func (l Link) Children() []codec.TLVInterface {
//     return [l.linkName, l.keyId, l.contentId]
// }

// String functions

func (l Link) String() string {
    return l.linkName.String()
}
