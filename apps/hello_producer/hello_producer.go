package main

import "fmt"

import "github.com/chris-wood/spud/stack/spud"
import "github.com/chris-wood/spud/stack/api/store"
import "github.com/chris-wood/spud/stack/api/portal"

func generateResponse(prefix string, payload []byte) []byte {
	fmt.Println("GENERATING A RESPONSE")
	return []byte("Hello, world!")
}

func serve(prefix string) {
	myStack, _ := spud.CreateRaw(`{"connector": "athena", "link": "tcp", "fwdaddress": "127.0.0.1:9696", "keys": ["key.p12"]}`)
	ccnPortal := portal.NewSecurePortal(myStack)
	api := store.NewStoreAPI(ccnPortal)

	done := make(chan int)

	api.Serve(prefix, generateResponse)

	<-done
}

func main() {
	serve("/hello")
}
