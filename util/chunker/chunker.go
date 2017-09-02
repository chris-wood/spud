package chunker

type Chunk []byte

type ChunkFunc func(acc interface{}, chunk Chunk) (interface{}, interface{})

type Chunker interface {
	GetChannel() chan Chunk
    Apply(f ChunkFunc, acc interface{}) interface{}
}
