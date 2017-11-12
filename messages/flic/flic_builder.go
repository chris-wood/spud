package flic

import (
	"github.com/chris-wood/spud/messages"
	"github.com/chris-wood/spud/messages/name"
	"github.com/chris-wood/spud/messages/hash"
	"github.com/chris-wood/spud/messages/flic/hashgroup"
	"github.com/chris-wood/spud/messages/payload"
	"github.com/chris-wood/spud/messages/content"
	"github.com/chris-wood/spud/codec"
	"github.com/chris-wood/spud/util/chunker"
	"crypto/sha256"
	"golang.org/x/crypto/hkdf"
	"io"
	"crypto/aes"
	"crypto/cipher"
	//"golang.org/x/crypto/xts"
)

// TODO(caw): the API should let clients specify a name, secret, and chunker -- CLEAN/FLIC handle the rest

func BuildEncryptedFLICTreeFromChunker(rootName name.Name, dataChunker chunker.Chunker) (*messages.MessageWrapper, []*messages.MessageWrapper) {
    digest := dataChunker.Hash(sha256.New())

    // Derive the CLEAN key
	info := []byte(rootName.String())
	salt := []byte{ } // TODO(caw): this would be the secret value provided by the producer, if available
	kdf := hkdf.New(sha256.New, digest, salt, info)

	masterKey := make([]byte, 32)
	_, err := io.ReadFull(kdf, masterKey)
	if err != nil {
		panic(err)
	}

	// Create the encryptor
	// TODO(caw): this should really be XTS
	block, err := aes.NewCipher(masterKey)
	if err != nil {
		panic(err)
	}
	iv := make([]byte, 12)
	encryptor := cipher.NewCTR(block, iv)


	// Create the encryptor function
	functor := func(encryptor interface{}, chunk chunker.Chunk) (interface{}, chunker.Chunk) {
		encryptedChunk := make(chunker.Chunk, len(chunk))
		encryptor.(cipher.Stream).XORKeyStream(encryptedChunk, chunk)
		return encryptor, encryptedChunk
	}
	newChunker := dataChunker.Map(functor, encryptor).(chunker.Chunker)

	return BuildFLICTreeFromChunker(newChunker)
}

// TODO(caw): this needs to accept a root name
func BuildFLICTreeFromChunker(dataChunker chunker.Chunker) (*messages.MessageWrapper, []*messages.MessageWrapper) {
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
