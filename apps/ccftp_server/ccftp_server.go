package main

import "os"
import "flag"
import "strings"

import "github.com/chris-wood/spud/stack/spud"
import "github.com/chris-wood/spud/stack/api/store"
import "github.com/chris-wood/spud/stack/api/portal"

type CCFTPServer struct {
	prefix string
}

func (s CCFTPServer) loadFile(name string, response []byte) []byte {
	if !strings.HasPrefix(name, s.prefix) {
		return make([]byte, 0)
	}

	file, err := os.Open("file.go") // For read access.
	if err != nil {
		return make([]byte, 0)
	}

	data := make([]byte, 100)
	count, err := file.Read(data)
	if err != nil {
		return make([]byte, 0)
	}

	err = file.Close()
	if err != nil {
		return make([]byte, 0)
	}

	return data[:count]
}

func (s CCFTPServer) serve(directory string) {
	myStack, _ := spud.CreateRaw("")
	ccnPortal := portal.NewPortal(myStack)
	api := store.NewStoreAPI(ccnPortal)
	api.Serve(s.prefix, s.loadFile)
}

func main() {
	baseDir := flag.String("dir", ".", "Path to the directory from which to serve files.")
	prefix := flag.String("prefix", "/ccftp/", "Producer server routable prefix.")

	flag.Parse()

	server := CCFTPServer{*prefix}
	server.serve(*baseDir)
}
