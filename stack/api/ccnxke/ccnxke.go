package ccnxke

import "fmt"
import "time"

import "github.com/chris-wood/spud/tables/lpm"
import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/util"
import "github.com/chris-wood/spud/stack/api/esic"
import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/kex"
import "github.com/chris-wood/spud/messages/interest"
import "github.com/chris-wood/spud/messages/content"

// XXX: wrap this in the crypto box
import "golang.org/x/crypto/nacl/box"

type CCNxKEAPI struct {
    kexStack stack.Stack
    prefixTable lpm.LPM
}

type SessionCallback func(session esic.ESIC)
type ResponseCallback func([]byte)

func NewCCNxKEAPI(s stack.Stack) *CCNxKEAPI {
    prefixLPM := lpm.LPM{}
    api := &CCNxKEAPI{
        kexStack: s,
        prefixTable: prefixLPM,
    }
    return api
}

func (api *CCNxKEAPI) Connect(prefix name.Name, handler SessionCallback) {
    fmt.Println("Starting the handshake...")

    randomPrefix, _ := util.GenerateRandomString(16)
    bareHelloName, _ := prefix.AppendComponent(randomPrefix)

    // Send the bare hello
    bareHello := kex.KEXHello()
    bareHelloRequest := interest.CreateWithName(bareHelloName)
    bareHelloRequest.AddContainer(bareHello)
    api.kexStack.Enqueue(bareHelloRequest)

    fmt.Println(bareHelloRequest)

    // Wait for the response, and use it to build the full hello
    time.Sleep(100 * time.Millisecond)
    reply := api.kexStack.Dequeue()

    fmt.Println("Got the reject, generating the full hello")

    reject, err := reply.GetContainer(codec.T_KEX)
    if err != nil {
        fmt.Println("Error: no KEX container in the REJECT content object")
        return
    }
    hello := kex.KEXFullHello(bareHello, reject.(*kex.KEX))

    randomPrefix, _ = util.GenerateRandomString(16)
    helloName, _ := prefix.AppendComponent(randomPrefix)

    helloRequest := interest.CreateWithName(helloName)
    helloRequest.AddContainer(hello)
    fmt.Println("Sending down the full hello and then sleeping...")
    api.kexStack.Enqueue(helloRequest)

    fmt.Println(helloRequest)

    // Wait for the response to complete the KEX
    time.Sleep(100 * time.Millisecond)
    reply = api.kexStack.Dequeue()

    fmt.Println("Got the accept, completing the handshake")

    accept, err := reply.GetContainer(codec.T_KEX)
    if err != nil {
        fmt.Println("Error: no KEX container in the ACCEPT content object")
        return
    }
    acceptKEX := accept.(*kex.KEX)
    fmt.Println(acceptKEX.GetContainerType())

    fmt.Println("Consumer: At KDF stage")

    var sharedKey [32]byte
    var peerPublic [32]byte
    var privateKey [32]byte
    copy(peerPublic[:], acceptKEX.GetPublicKeyShare())
    copy(privateKey[:], hello.GetPrivateKeyShare())
    box.Precompute(&sharedKey, &peerPublic, &privateKey)

    fmt.Println(sharedKey)
}

func (api *CCNxKEAPI) Service(prefix name.Name, callback SessionCallback) {
    macKey, _ := util.GenerateRandomBytes(16)
    encKey, _ := util.GenerateRandomBytes(16)
    go api.serviceSessions(prefix, callback, macKey, encKey)
}

func (api *CCNxKEAPI) serviceSessions(prefix name.Name, callback SessionCallback, macKey, encKey []byte) {
    for ;; {
        // Wait for a hello
        request := api.kexStack.Dequeue()
        if !prefix.IsPrefix(request.Name()) {
            break
        }

        // Handle this message
        kexContainer, err := request.GetContainer(codec.T_KEX)
        

        // bareHello, err := request.GetContainer(codec.T_KEX)
        // if err != nil {
        //     fmt.Println("Error: no KEX container in the BARE HELLO interest")
        //     return
        // }
        // reject := kex.KEXHelloReject(bareHello.(*kex.KEX), macKey)
        // rejectResponse := content.CreateWithName(request.Name())
        // rejectResponse.AddContainer(reject)
        // api.kexStack.Enqueue(rejectResponse)
        //
        // fmt.Println("sending down the rejection message and then sleeping")
        //
        // // Wait for the full hello to come back
        // time.Sleep(100 * time.Millisecond)
        // request = api.kexStack.Dequeue()
        //
        // fmt.Println("Got the full hello, completing the handshake")
        //
        // hello, err := request.GetContainer(codec.T_KEX)
        // if err != nil {
        //     fmt.Println("Error: no KEX container in the HELLO interest")
        //     return
        // }
        // fmt.Println("generating the accept")
        // accept := kex.KEXHelloAccept(bareHello.(*kex.KEX), reject, hello.(*kex.KEX), macKey, encKey)
        // if accept == nil {
        //     // XXX: implement better recovery here
        //     fmt.Println("recover")
        // }
        // acceptResponse := content.CreateWithName(request.Name())
        // acceptResponse.AddContainer(accept)
        //
        // api.kexStack.Enqueue(acceptResponse)
        //
        // // Create the session and give it to the callback
        // // XXX
        // fmt.Println("Producer: At KDF stage")
        // var sharedKey [32]byte
        // var peerPublic [32]byte
        // var privateKey [32]byte
        // helloKex := hello.(*kex.KEX)
        // copy(peerPublic[:], helloKex.GetPublicKeyShare())
        // copy(privateKey[:], accept.GetPrivateKeyShare())
        // box.Precompute(&sharedKey, &peerPublic, &privateKey)
        //
        // fmt.Println(sharedKey)


    }
}
