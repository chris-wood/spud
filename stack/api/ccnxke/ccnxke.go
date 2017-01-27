package ccnxke

import "log"
import "time"

import "github.com/chris-wood/spud/tables/lpm"
import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/util"
import "github.com/chris-wood/spud/stack/api/esic"
import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/kex"
import "github.com/chris-wood/spud/messages/interest"
import "github.com/chris-wood/spud/messages/content"

// XXX: wrap this in the crypto box
import "golang.org/x/crypto/nacl/box"

type CCNxKEAPI struct {
    kexStack *stack.Stack
    prefixTable lpm.LPM
}

type SessionCallback func(session *esic.ESIC)
type ResponseCallback func([]byte)

const connectString string = "CONNECT"

func NewCCNxKEAPI(s *stack.Stack) *CCNxKEAPI {
    prefixLPM := lpm.LPM{}
    api := &CCNxKEAPI{
        kexStack: s,
        prefixTable: prefixLPM,
    }
    return api
}

func (api *CCNxKEAPI) Connect(prefix name.Name, handler SessionCallback) {
    randomPrefix, _ := util.GenerateRandomString(16)
    bareHelloName, _ := prefix.AppendComponent(connectString)
    bareHelloName, _ = bareHelloName.AppendComponent(randomPrefix)

    // Send the bare hello
    bareHello := kex.KEXHello()
    bareHelloRequest := interest.CreateWithName(bareHelloName)
    bareHelloRequest.AddContainer(bareHello)
    api.kexStack.Enqueue(messages.Package(bareHelloRequest))

    // Wait for the response, and use it to build the full hello
    time.Sleep(100 * time.Millisecond)
    replyWrapper := api.kexStack.Dequeue()
    reply := replyWrapper.InnerMessage()
    log.Println("Got the REJECT")

    reject, err := reply.GetContainer(codec.T_KEX)
    if err != nil {
        log.Println("Error: no KEX container in the REJECT content object")
        return
    }
    hello := kex.KEXFullHello(bareHello, reject.(*kex.KEX))

    randomPrefix, _ = util.GenerateRandomString(16)
    helloName, _ := prefix.AppendComponent(connectString)
    helloName, _ = prefix.AppendComponent(randomPrefix)

    helloRequest := interest.CreateWithName(helloName)
    helloRequest.AddContainer(hello)
    api.kexStack.Enqueue(messages.Package(helloRequest))

    // Wait for the response to complete the KEX
    time.Sleep(100 * time.Millisecond)
    replyWrapper = api.kexStack.Dequeue()
    reply = replyWrapper.InnerMessage()

    accept, err := reply.GetContainer(codec.T_KEX)
    if err != nil {
        log.Println("Error: no KEX container in the ACCEPT content object")
        return
    }
    acceptKEX := accept.(*kex.KEX)

    var sharedKey [32]byte
    var peerPublic [32]byte
    var privateKey [32]byte
    copy(peerPublic[:], acceptKEX.GetPublicKeyShare())
    copy(privateKey[:], hello.GetPrivateKeyShare())
    box.Precompute(&sharedKey, &peerPublic, &privateKey)

    // Create and start the session
    // session := esic.NewESIC(api.kexStack, sharedKey[:], acceptKEX.GetSessionID())
    // handler(session)

    log.Println("Consumer: ", sharedKey)
}

func (api *CCNxKEAPI) Service(prefix name.Name, callback SessionCallback) {
    macKey, _ := util.GenerateRandomBytes(16)
    encKey, _ := util.GenerateRandomBytes(16)
    go api.serviceSessions(prefix, callback, macKey, encKey)
}

func (api *CCNxKEAPI) serviceSessions(prefix name.Name, callback SessionCallback, macKey, encKey []byte) {
    for ;; {
        requestWrapper := api.kexStack.Dequeue()
        request := requestWrapper.InnerMessage()
        if !prefix.IsPrefix(request.Name()) {
            break
        }

        kexTLV, _ := request.GetContainer(codec.T_KEX)
        kexContainer := kexTLV.(*kex.KEX)

        switch kexContainer.GetMessageType() {
        case codec.T_KEX_BAREHELLO:
            log.Println("Got the BARE HELLO")
            reject := kex.KEXHelloReject(kexContainer, macKey)
            rejectResponse := content.CreateWithName(request.Name())
            rejectResponse.AddContainer(reject)
            api.kexStack.Enqueue(messages.Package(rejectResponse))
            break

        case codec.T_KEX_HELLO:
            log.Println("Got the HELLO")
            accept, err := kex.KEXHelloAccept(kexContainer, macKey, encKey)
            if err != nil {
                log.Println(err)
                break
            }

            acceptResponse := content.CreateWithName(request.Name())
            acceptResponse.AddContainer(accept)
            api.kexStack.Enqueue(messages.Package(acceptResponse))

            // XXX: go to the KDF step

            var sharedKey [32]byte
            var peerPublic [32]byte
            var privateKey [32]byte
            copy(peerPublic[:], kexContainer.GetPublicKeyShare())
            copy(privateKey[:], accept.GetPrivateKeyShare())
            box.Precompute(&sharedKey, &peerPublic, &privateKey)

            log.Println("Producer:", sharedKey)

            // Create and start the session
            // session := esic.NewESIC(api.kexStack, sharedKey[:], accept.GetSessionID())
            // callback(session)

            break

        case codec.T_KEX_REJECT:
        case codec.T_KEX_ACCEPT:
            log.Println("Got an invalid message...")
            // invalid message type to be received here...
            break
        }

        time.Sleep(100 * time.Millisecond)
    }
}
