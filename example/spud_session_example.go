package main

import (
	"github.com/chris-wood/spud/stack/api/portal"
	"github.com/chris-wood/spud/messages/name"
	"github.com/chris-wood/spud/stack/spud"
	"github.com/chris-wood/spud/messages"
	"github.com/chris-wood/spud/messages/interest"
	"fmt"
	"github.com/chris-wood/spud/messages/content"
	"github.com/chris-wood/spud/messages/payload"
)

var done chan int

func displayResponse(response *messages.MessageWrapper) {
	fmt.Println("Response: " + string(response.Payload().Value()))
	done <- 1
}

func generateResponse(request *messages.MessageWrapper) *messages.MessageWrapper {
	data := []byte("Hello world")
	dataPayload := payload.Create(data)
	response := messages.Package(content.CreateWithNameAndPayload(request.Name(), dataPayload))
	return response
}

 func testSession() {
	 myStack, err := spud.CreateRaw(`{"connector": "athena", "link": "loopback", "fwd-address": "127.0.0.1:9696", "keys": ["key.p12"]}`)
	 if err != nil {
		 panic("Could not create the stack")
	 }
	 done = make(chan int)

	 p := portal.NewSecurePortal(myStack)

	 prefix, _ := name.Parse("ccnx:/producer")

	 go p.Serve(prefix, generateResponse)
	 p.Connect(prefix)

	 fmt.Println("SENDING REQUEST FOR", prefix)
	 requestInterest := interest.CreateWithName(prefix)
	 requestWrapper := messages.Package(requestInterest)
	 p.GetAsync(requestWrapper, displayResponse)

	 // sleep until the consumer gets a response
	 <- done
 }

func main() {
	testSession()
}
