package flic

import "github.com/chris-wood/spud/util/chunker"
import "io/ioutil"

import "testing"

func TestCreate(t *testing.T) {
	dataSize := 1024
	data := make([]byte, dataSize)
	for i := 0; i < dataSize; i++ {
		data[i] = uint8(i)
	}

	chunkSize := 256
	fname := "/tmp/flic_data_test"

	err := ioutil.WriteFile(fname, data, 0644)
	if err != nil {
		t.Errorf("Failed to write data to the file")
	}

	fChunker, err := chunker.NewFileChunker(fname, chunkSize)
	if err != nil {
		t.Error("Unable to create the file chunker:", err)
	}

	root, messages := CreateFLICTreeFromChunker(fChunker)

    if root == nil {
        t.Error("Root is invalid")
    }

    expected := 5
	if len(messages) != expected {
		t.Error("Invalid message collection returned, got {}, expected {}", expected, len(messages))
	}
	t.Log(len(messages))
}

func TestLookup(t *testing.T) {
	// pass
}
