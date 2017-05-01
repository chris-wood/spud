package main

import (
	"github.com/chris-wood/spud/stack/api/portal"
	"github.com/chris-wood/spud/messages/name"
	"github.com/chris-wood/spud/stack/spud"
	"github.com/chris-wood/spud/messages"
	"github.com/chris-wood/spud/messages/interest"
	"fmt"
)

var done chan int

func displayResponse(response *messages.MessageWrapper) {
	fmt.Println("Response: " + string(response.Payload().Value()))
	done <- 1
}

func generateResponse(request *messages.MessageWrapper) *messages.MessageWrapper {
	return nil
}

 func testSession() {
	 myStack, err := spud.CreateRaw(`{"connector": "athena", "link": "loopback", "fwd-address": "127.0.0.1:9696", "keys": ["key.p12"]}`)
	 if err != nil {
		 panic("Could not create the stack")
	 }
	 done = make(chan int)

	 p := portal.NewSecurePortal(myStack)

	 prefix, _ := name.Parse("ccnx:/producer")

	 p.Serve(prefix, generateResponse) // ditto below

	 p.Connect(prefix)

	 requestInterest := interest.CreateWithName(prefix)
	 requestWrapper := messages.Package(requestInterest)

	 p.GetAsync(requestWrapper, displayResponse)

	 // sleep until the consumer gets a response
	 <- done
 }

func main() {
	testSession()
}
