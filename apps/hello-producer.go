package main

import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/stack/api/adapter"

func generateResponse(prefix string, payload []byte) []byte {
    return []byte("Hello, world!")
}

func serve(prefix string) {
    myStack := stack.Create("")
    api := adapter.NewNameAPI(myStack)

    done := make(chan int)

    api.Serve(prefix, generateResponse)

    <-done
}

func main() {
    serve("/hello")
}
