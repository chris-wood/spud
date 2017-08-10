package filter

import "github.com/chris-wood/spud/util/chunker"

type Filter interface {
	Apply(c chunker.Chunker) chunker.Chunker
}
