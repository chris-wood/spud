package main

import "fmt"
import "time"

import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/stack/api/adapter"
import "github.com/chris-wood/spud/stack/api/ccnxke"
import "github.com/chris-wood/spud/stack/api/esic"

import "github.com/chris-wood/spud/messages/name"

var count = 0
var done chan int

func displayResponse(response []byte) {
    fmt.Println("Response: " + string(response))
    count = 1
}

func generateResponse(name string, response []byte) []byte {
    return []byte("hello, spud!")
}

func testStack() {
    // myStack := stack.Create(`{"link": "loopback"}`)
    myStack := stack.Create(`{"connector": "athena", "link": "loopback", "fwd-address": "127.0.0.1:9696", "keys": ["key.p12"]}`)
    api := adapter.NewNameAPI(myStack)

    api.Serve("ccnx:/hello/spud", generateResponse)
    api.Get("ccnx:/hello/spud", displayResponse)

    for ;; {
        if count == 0 {
            time.Sleep(100 * time.Millisecond)
        } else {
            break
        }
    }
}


func ProducerSessionHandler(session *esic.ESIC) {
    session.Serve("/foo/bar", func(nameString string, data []byte) []byte {
        return []byte("Hello CCNxKE!")
    })
}

func ConsumerSessionHandler(session *esic.ESIC) {
    session.Get("/foo/bar", func(data []byte) {
        done <- 1
        fmt.Println("Received:", string(data))
    })
}

func testSession() {
    myStack := stack.Create(`{"link": "loopback"}`)
    api := ccnxke.NewCCNxKEAPI(myStack)
    done = make(chan int)

    prefix, _ := name.Parse("ccnx:/producer")
    api.Service(prefix, ProducerSessionHandler) // ditto below
    api.Connect(prefix, ConsumerSessionHandler) // SessionHandler will be invoked if and when the session is successfully completed

    // sleep until the consumer gets a response
    <- done
}

func main() {
    testStack()
    // testSession()
}
