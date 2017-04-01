package chunker

type Chunk []byte

type Chunker interface {
	GetChannel() chan Chunk
}
