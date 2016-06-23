package payload

// import "github.com/chris-wood/spud/codec"
import "fmt"
// import "encoding/json"

type Payload struct {
    bytes []byte
}

type payloadError struct {
    prob string
}

func (e payloadError) Error() string {
    return fmt.Sprintf("%s", e.prob)
}

// Constructors

func Create(bytes []byte) *Payload {
    return &Payload{payloadType, bytes}
}

// TLV functions

func (p Payload) Type() uint16 {
    return uint16(codec.T_PAYLOAD)
}

func (p Payload) TypeString() string {
    return "Payload"
}

func (p Payload) Length() uint16 {
    return len(bytes)
}

func (p Payload) Value() []byte  {
    return p.bytes
}

func (p Payload) Children() []codec.TLV {
    return nil
}

// String functions

func (p Payload) String() string {
    return string(p.bytes)
}
