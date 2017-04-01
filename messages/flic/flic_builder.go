package flic

import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/messages/hash"
import "github.com/chris-wood/spud/messages/content"
import "github.com/chris-wood/spud/util/chunker"

import "crypto/sha256"

func CreateFLICTreeFromChunker(dataChunker chunker.Chunker) []*messages.MessageWrapper {
	root := CreateEmptyHashGroup()
	collection := make([]*messages.MessageWrapper, 0)

	dataSize := 0
	for chunk := range dataChunker.GetChannel() {
		dataPayload := payload.Create(chunk)
		namelessLeaf := content.CreateWithPayload(dataPayload)
		namelessMsg := messages.Package(namelessLeaf)
		collection = append(collection, namelessMsg)

		leafHashRaw := namelessMsg.ComputeMessageHash(sha256.New())
		leafHash := hash.Create(hash.HashTypeSHA256, leafHashRaw)
		leafPointer := SizedDataPointer{size: Size{codec.T_POINTER_SIZE, uint64(len(chunk))}, ptrHash: leafHash}

		dataSize = dataSize + len(chunk)

		if ok := root.AddPointer(leafPointer); !ok {
			parentFlic := CreateFLICFromHashGroup(root)
			parentMsg := messages.Package(parentFlic)
			parentHashRaw := parentMsg.ComputeMessageHash(sha256.New())
			parentHash := hash.Create(hash.HashTypeSHA256, parentHashRaw)
			parentPointer := SizedManifestPointer{size: Size{codec.T_POINTER_SIZE, uint64(dataSize)}, ptrHash: parentHash}

			// Reset the data size
			dataSize = 0

			newRoot := CreateEmptyHashGroup()
			newRoot.AddPointer(parentPointer)

			root = newRoot
		}
	}

	rootFLIC := CreateFLICFromHashGroup(root)
	rootMessage := messages.Package(rootFLIC)
	collection = append(collection, rootMessage)

	return collection
}
