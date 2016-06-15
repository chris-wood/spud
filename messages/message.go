package messages

type Message interface {
    Type() int
    Version() int
    Length() int

    Encode() []byte
}
