package main

import "flag"
import "fmt"

import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/stack/api/adapter"

type CCFTPFetcher struct {
    prefix string
}

func displayResponse(response []byte) {
    fmt.Println("Response: " + string(response))
}

func (f CCFTPFetcher) fetch(file string) {
    myStack := stack.Create("")
    api := adapter.NewNameAPI(myStack)

    // XXX: build the name based on the prefix and file

    api.Get("ccnx:/hello/spud", displayResponse)
}

func main() {
    fileName := flag.String("file", ".", "Name of the file to fetch.")
    prefix := flag.String("prefix", "/ccftp/", "Producer server routable prefix.")
    flag.Parse()

    fetcher := CCFTPFetcher{*prefix}
    fetcher.fetch(*fileName)
}
