package main

import "fmt"
import "bufio"
import "os"

import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/stack/api/adapter"

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

func main() {
    testStack()
}
