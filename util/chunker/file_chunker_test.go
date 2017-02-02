package chunker

import "testing"
import "io/ioutil"

func TestFileChunker(t *testing.T) {
    data := make([]byte, 1024)
    for i := 0; i < 1024; i++ {
        data[i] = uint8(i)
    }

    chunkSize := 256 // sizeof(uint8)
    fname := "/tmp/file_chunker_test"

    err := ioutil.WriteFile(fname, data, 0644)
    if err != nil {
        t.Errorf("Failed to write data to the file")
    }


    fChunker, err := NewFileChunker(fname, chunkSize)
    if err != nil {
        t.Error("Unable to create the file chunker:", err)
    }

    channel := fChunker.GetChannel()
    for chunk := range(channel) {
        if len(chunk) != chunkSize {
            t.Errorf("Incorrect chunk size. Got %d, expected %d", len(chunk), chunkSize)
        }
        for i := 0; i < chunkSize; i++ {
            if chunk[i] != uint8(i) {
                t.Errorf("Incorrect chunk data. Got %d at index %d, expected %d", int(chunk[i]), i, i)
            }
        }
    }
}
