package chunker

import "hash"

type Chunk []byte

type Chunker interface {
	GetChannel() chan Chunk
    // Hash(hasher hash.Hash) []byte
}
