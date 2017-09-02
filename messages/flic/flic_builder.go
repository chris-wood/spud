package flic

import (
	"github.com/chris-wood/spud/messages"
	"github.com/chris-wood/spud/messages/hash"
	"github.com/chris-wood/spud/messages/flic/hashgroup"
	"github.com/chris-wood/spud/messages/payload"
	"github.com/chris-wood/spud/messages/content"
	"github.com/chris-wood/spud/codec"
	"github.com/chris-wood/spud/util/chunker"
	"crypto/sha256"
)

func CreateFLICTreeFromChunker(dataChunker chunker.Chunker) (*messages.MessageWrapper, []*messages.MessageWrapper) {
	root := hashgroup.CreateEmptyHashGroup()
	collection := make([]*messages.MessageWrapper, 0)

	dataSize := 0
	for chunk := range dataChunker.GetChannel() {
		dataPayload := payload.Create(chunk)
		namelessLeaf := content.CreateWithPayload(dataPayload)
		namelessMsg := messages.Package(namelessLeaf)
		collection = append(collection, namelessMsg)

		leafHashRaw := namelessMsg.ComputeMessageHash(sha256.New())
		leafHash := hash.Create(hash.HashTypeSHA256, leafHashRaw)
		leafPointer := hashgroup.CreateSizedDataPointer(hashgroup.CreateSize(codec.T_POINTER_SIZE, uint64(len(chunk))), leafHash)

		dataSize = dataSize + len(chunk)

		if ok := root.AddPointer(leafPointer); !ok {
			parentFlic := CreateFLICFromHashGroup(root)
			parentMsg := messages.Package(parentFlic)
			parentHashRaw := parentMsg.ComputeMessageHash(sha256.New())
			parentHash := hash.Create(hash.HashTypeSHA256, parentHashRaw)
			parentPointer := hashgroup.CreateSizedManifestPointer(hashgroup.CreateSize(codec.T_POINTER_SIZE, uint64(dataSize)), parentHash)

			// Reset the data size
			dataSize = 0

			newRoot := hashgroup.CreateEmptyHashGroup()
			newRoot.AddPointer(parentPointer)

			root = newRoot
		}
	}

	rootFLIC := CreateFLICFromHashGroup(root)
	rootMessage := messages.Package(rootFLIC)
	collection = append(collection, rootMessage)

	return rootMessage, collection
}
