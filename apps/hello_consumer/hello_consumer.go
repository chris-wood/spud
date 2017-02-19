package main

import "fmt"

import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/stack/spud"
import "github.com/chris-wood/spud/stack/api/kvs"
import "github.com/chris-wood/spud/stack/api/portal"

var done chan int

func displayResponse(response []byte) {
    fmt.Println("Response: " + string(response))
    done <- 1
}

func get(nameString string) {
    myStack, _ := spud.CreateRaw(`{"connector": "athena", "link": "tcp", "fwdaddress": "127.0.0.1:9695", "keys": ["key.p12"]}`)
    ccnPortal := portal.NewSecurePortal(myStack)
    prefix, _ := name.Parse(nameString)
    ccnPortal.Connect(prefix)
    api := adapter.NewKVSAPI(ccnPortal)

    done = make(chan int)

    fmt.Println("Fetching now...")
    api.GetAsync(nameString, displayResponse)

    <- done
}

func main() {
    get("/hello")
}
