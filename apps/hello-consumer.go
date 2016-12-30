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
    myStack := stack.Create("")
    api := adapter.NewNameAPI(myStack)

    done = make(chan int)

    api.Get(name, displayResponse)

    <- done
}

func main() {
    get("/hello")
}
