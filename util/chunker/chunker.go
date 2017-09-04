package chunker

import "hash"

type Chunk []byte

type ChunkMapFunc func(acc interface{}, chunk Chunk) (interface{}, Chunk)
type ChunkFunc func(acc interface{}, chunk Chunk) (interface{}, interface{})

type Chunker interface {
	GetChannel() chan Chunk
	Map(f ChunkMapFunc, acc interface{}) Chunker
    Apply(f ChunkFunc, acc interface{}) interface{}
    Hash(hasher hash.Hash) []byte
}
