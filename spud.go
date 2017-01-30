package main

import "fmt"
import "time"

import "github.com/chris-wood/spud/stack/spud"
import "github.com/chris-wood/spud/stack/api/kvs"
import "github.com/chris-wood/spud/stack/api/portal"

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
    myStack, err := spud.CreateRaw(`{"connector": "athena", "link": "loopback", "fwd-address": "127.0.0.1:9696", "keys": ["key.p12"]}`)
    if err != nil {
        panic("Could not create the stack")
    }

    p := portal.NewPortal(myStack)
    api := adapter.NewKVSAPI(p)

    api.Serve("ccnx:/hello/spud", generateResponse)
    data, err := api.Get("ccnx:/hello/spud", time.Second)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(string(data))
    }
}

// func testSession() {
//     myStack, err := spud.CreateRaw(`{"link": "loopback"}`)
//     if err != nil {
//         panic("Could not create the stack")
//     }
//     done = make(chan int)
//
//     prefix, _ := name.Parse("ccnx:/producer")
//     api.Service(prefix) // ditto below
//     api.Connect(prefix) // SessionHandler will be invoked if and when the session is successfully completed
//
//     // sleep until the consumer gets a response
//     <- done
// }

func main() {
    testStack()
    // testSession()
}
