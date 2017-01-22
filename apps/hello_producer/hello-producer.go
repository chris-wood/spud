package main

import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/stack/api/kvs"
import "github.com/chris-wood/spud/stack/api/portal"

func generateResponse(prefix string, payload []byte) []byte {
    return []byte("Hello, world!")
}

func serve(prefix string) {
    myStack, _ := stack.CreateRaw(`{"connector": "athena", "link": "tcp", "fwd-address": "127.0.0.1:9696", "keys": ["key.p12"]}`)
    // myStack := stack.Create(`{"connector": "athena", "link": "loopback", "fwd-address": "127.0.0.1:9696", "keys": ["key.p12"]}`)
    // myStack := stack.CreateTest()
    ccnPortal := portal.NewPortal(myStack)
    api := adapter.NewKVSAPI(ccnPortal)

    done := make(chan int)

    api.Serve(prefix, generateResponse)

    <-done
}

func main() {
    serve("/hello")
}
