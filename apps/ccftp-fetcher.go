package main

import "flag"
import "fmt"

import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/stack/api/kvs"
import "github.com/chris-wood/spud/stack/api/portal"

type CCFTPFetcher struct {
    prefix string
}

func displayResponse(response []byte) {
    fmt.Println("Response: " + string(response))
}

func (f CCFTPFetcher) fetch(file string) {
    myStack, _ := stack.CreateRaw("")
    ccnPortal := portal.NewPortal(myStack)
    api := adapter.NewKVSAPI(ccnPortal)

    // XXX: build the name based on the prefix and file

    api.GetAsync("ccnx:/hello/spud", displayResponse)
}

func main() {
    fileName := flag.String("file", ".", "Name of the file to fetch.")
    prefix := flag.String("prefix", "/ccftp/", "Producer server routable prefix.")
    flag.Parse()

    fetcher := CCFTPFetcher{*prefix}
    fetcher.fetch(*fileName)
}
