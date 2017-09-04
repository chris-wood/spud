package chunker

import "os"
import "bufio"
import "hash"

type FileChunker struct {
	Fname        string
	Size         int
	NumChunks    int
	chunkChannel chan Chunk
}

func NewFileChunker(fname string, chunkSize int) (*FileChunker, error) {
	fh, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	fi, err := fh.Stat()
	if err != nil {
		return nil, err
	}

	totalSize := int(fi.Size())
	numChunks := ((totalSize - 1) / chunkSize) + 1
	chunks := make(chan Chunk, numChunks)

	reader := bufio.NewReader(fh)
	buffer := make([]byte, chunkSize)

	for i := 0; i < numChunks; i++ {
		numRead, err := reader.Read(buffer)
		if err != nil {
			return nil, err
		}
		slice := buffer[0:numRead]
		chunks <- Chunk(slice)
	}

	// If we got this far, close up the channel
	close(chunks)

	return &FileChunker{fname, chunkSize, numChunks, chunks}, nil
}

func (f *FileChunker) GetChannel() chan Chunk {
	return f.chunkChannel
}

func (f *FileChunker) Apply(chunkFunc ChunkFunc, acc interface{}) interface{} {
    tmp := make(chan Chunk, f.NumChunks)
    var result interface{}
    for chunk := range f.chunkChannel {
        acc, result = chunkFunc(acc, chunk)
        tmp <- chunk // Append back to the normal channel
    }
    close(tmp)

    // Reset the chunk channel
    f.chunkChannel = tmp

    return result
}

func (f *FileChunker) Map(chunkFunc ChunkMapFunc, acc interface{}) Chunker {
	tmp := make(chan Chunk, f.NumChunks)
	var newChunk Chunk
	for chunk := range f.chunkChannel {
		acc, newChunk = chunkFunc(acc, chunk)
		tmp <- newChunk // Append back to the normal channel
	}
	close(tmp)

	return &FileChunker{f.Fname, f.Size, f.NumChunks, tmp}
}

func (f *FileChunker) Hash(hasher hash.Hash) []byte {
    tmp := make(chan Chunk, f.NumChunks)
    for chunk := range f.chunkChannel {
        hasher.Write(chunk)
        tmp <- chunk // Append back to the normal channel
    }
    f.chunkChannel = tmp

    return hasher.Sum(nil)
}
