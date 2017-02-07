package flic

import "github.com/chris-wood/spud/util/chunker"
import "io/ioutil"

import "testing"

func TestCreate(t *testing.T) {
    data := make([]byte, 1024)
    for i := 0; i < 1024; i++ {
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

    messages := CreateFLICTreeFromChunker(fChunker)

    // XXX: compute the expected number of FLIC entries

    if len(messages) == 0 {
        t.Error("Invalid message collection returned")
    }
    t.Log(len(messages))
}

func TestLookup(t *testing.T) {
    // pass
}
