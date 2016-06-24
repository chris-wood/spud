package adapter

// TODO: names+options in, data out
type Adapter interface {
    Get(name string)
    Put(data []byte)
}
