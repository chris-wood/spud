package spud

import "fmt"
import "time"

import "github.com/chris-wood/spud/stack/spud"
import "github.com/chris-wood/spud/stack/api/store"
import "github.com/chris-wood/spud/stack/api/portal"

var count = 0
var done chan int

func displayResponse(response []byte) {
	fmt.Println("Response: " + string(response))
	count = 1
}

func generateResponse(name string, response []byte) []byte {
	fmt.Println("here's the response")
	return []byte("hello, spud!")
}

func testStack() {
	myStack, err := spud.CreateRaw(`{"connector": "athena", "link": "loopback", "fwd-address": "127.0.0.1:9696", "keys": ["key.p12"]}`)
	if err != nil {
		panic("Could not create the stack")
	}

	p := portal.NewPortal(myStack)
	storer := store.NewStoreAPI(p)

	storer.Serve("ccnx:/hello/spud", generateResponse)
	data, err := storer.Get("ccnx:/hello/spud", time.Second)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(data))
	}
}

func main() {
	testStack()
}
