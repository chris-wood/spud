package main

import "fmt"
import "bufio"
import "os"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/messages/interest"
import "github.com/chris-wood/spud/messages/content"
import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/stack/adapter"

import "crypto/rand"

func simpleTest() {
    name1, err := name.Parse("ccnx:/hello/spud")
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(name1.Length());

    e := codec.Encoder{}
    nameBytes := e.EncodeTLV(name1);
    fmt.Println(nameBytes)

    d := codec.Decoder{}

    nameTLV := d.Decode(nameBytes)
    decodedName, err := name.CreateFromTLV(nameTLV[0])
    fmt.Println(decodedName)

    interestMessage := interest.CreateWithName(decodedName)
    interestBytes := e.EncodeTLV(interestMessage)
    fmt.Println(interestBytes)
}

func displayResponse(response []byte) {
    fmt.Println("Response: " + string(response))
}

func generateResponse(name string, response []byte) []byte {
    b := make([]byte, 32)
    rand.Read(b)
    return b
}

func testStack() {
    myStack := stack.Create("")
    api := adapter.NewNameAPI(myStack)

    api.Serve("/hello/spud/", generateResponse)
    api.Get("ccnx:/hello/spud", displayResponse)

    reader := bufio.NewReader(os.Stdin)
    for ;; {
        fmt.Print("> ")
        text, _ := reader.ReadString('\n')
        fmt.Println(text)
    }
}

func main() {
    name1, err := name.Parse("ccnx:/hello/spud")
    if err != nil {
        fmt.Println(err)
    }

    bytes := make([]byte, 32)
    rand.Read(bytes)
    dataPayload := payload.Create(bytes)

    contentMsg := content.CreateWithNameAndPayload(name1, dataPayload)
    fmt.Println("Content: " + contentMsg.Identifier())

    e := codec.Encoder{}
    contentBytes := e.EncodeTLV(contentMsg);

    fmt.Println(contentBytes)
}
