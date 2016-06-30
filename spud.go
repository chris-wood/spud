package main

import "fmt"
import "bufio"
import "os"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/interest"
import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/stack/adapter"

func displayResponse(response []byte) {
    fmt.Println(string(response))
}

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

func main() {
    myStack := stack.Create("")
    api := adapter.NewNameAPI(myStack)
    api.Get("ccnx:/hello/spud", displayResponse)

    reader := bufio.NewReader(os.Stdin)
    for ;; {
        fmt.Print("> ")
        text, _ := reader.ReadString('\n')
        fmt.Println(text)
    }
}
