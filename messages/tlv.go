package messages

type TLV interface {
    Type() int
    Length() int
    Value() []byte
}
