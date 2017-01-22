package main

import "fmt"

import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/stack/api/kvs"
import "github.com/chris-wood/spud/stack/api/portal"

var done chan int

func displayResponse(response []byte) {
    fmt.Println("Response: " + string(response))
    done <- 1
}

func get(name string) {
    myStack, _ := stack.CreateRaw(`{"connector": "athena", "link": "tcp", "fwd-address": "127.0.0.1:9695", "keys": ["key.p12"]}`)
    // myStack := stack.Create(`{"connector": "athena", "link": "loopback", "fwd-address": "127.0.0.1:9696", "keys": ["key.p12"]}`)
    // myStack := stack.CreateTest()
    ccnPortal := portal.NewPortal(myStack)
    api := adapter.NewKVSAPI(ccnPortal)

    done = make(chan int)

    fmt.Println("Fetching now...")
    api.GetAsync(name, displayResponse)

    <- done
}

func main() {
    get("/hello")
}
