package main

import "fmt"

import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/stack/api/adapter"

var done chan int

func displayResponse(response []byte) {
    fmt.Println("Response: " + string(response))
    done <- 1
}

func get(name string) {
    myStack := stack.Create(`{"connector": "athena", "link": "tcp", "fwd-address": "127.0.0.1:9695", "keys": ["key.p12"]}`)
    // myStack := stack.CreateTest()
    api := adapter.NewNameAPI(myStack)

    done = make(chan int)

    fmt.Println("Fetching now...")
    api.Get(name, displayResponse)

    <- done
}

func main() {
    get("/hello")
}
