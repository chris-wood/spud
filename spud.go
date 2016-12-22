package main

import "fmt"
import "bufio"
import "time"
import "os"

import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/stack/api/adapter"
import "github.com/chris-wood/spud/stack/api/ccnxke"
import "github.com/chris-wood/spud/stack/api/esic"

import "github.com/chris-wood/spud/messages/name"

func displayResponse(response []byte) {
    fmt.Println("Response: " + string(response))
}

func generateResponse(name string, response []byte) []byte {
    return []byte("hello, spud!")
}

func testStack() {
    myStack := stack.Create("")
    api := adapter.NewNameAPI(myStack)

    api.Serve("ccnx:/hello/spud", generateResponse)
    api.Get("ccnx:/hello/spud", displayResponse)

    reader := bufio.NewReader(os.Stdin)
    for ;; {
        fmt.Print("> ")
        reader.ReadString('\n') // text, err :=
    }
}

var count = 0

func ProducerSessionHandler(session *esic.ESIC) {
    fmt.Println("Producer session established! ")
    count++

    session.Serve("/foo/bar", func(nameString string, data []byte) []byte {
        fmt.Println("Producer supplying a response...")
        return []byte("Hello CCNxKE!")
    })
}

var done chan int

func ConsumerSessionHandler(session *esic.ESIC) {
    fmt.Println("Consumer session established! ")
    count++

    session.Get("/foo/bar", func(data []byte) {
        fmt.Println("Got the data back:", string(data))
        done <- 1
    })
}

func testSession() {
    myStack := stack.Create("")
    api := ccnxke.NewCCNxKEAPI(myStack)

    prefix, _ := name.Parse("ccnx:/producer")
    api.Service(prefix, ProducerSessionHandler) // ditto below
    api.Connect(prefix, ConsumerSessionHandler) // SessionHandler will be invoked if and when the session is successfully completed

    for ;; {
        if count == 0 {
            time.Sleep(500 * time.Millisecond)
        } else {
            break
        }
    }

    // sleep until the consumer gets a response
    <- done
}

func main() {
    // testStack()
    testSession()
}
