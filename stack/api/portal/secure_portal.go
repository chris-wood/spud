package portal

import "time"
import "log"

import "github.com/chris-wood/spud/stack"
import "github.com/chris-wood/spud/util/random"
import "github.com/chris-wood/spud/stack/component/tunnel"
import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/kex"
import "github.com/chris-wood/spud/messages/interest"
import "github.com/chris-wood/spud/messages/content"

// XXX: wrap this in the crypto box
import "golang.org/x/crypto/nacl/box"

const connectString string = "CONNECT"

type SecurePortal struct {
	apiStack stack.Stack
	macKey   []byte
	encKey   []byte
}

func NewSecurePortal(s stack.Stack) SecurePortal {
	api := SecurePortal{
		apiStack: s,
	}

	return api
}

func (n SecurePortal) Connect(prefix *name.Name) {
	randomSuffix, _ := random.GenerateRandomString(16)
	bareHelloName, _ := prefix.AppendComponent(connectString)
	bareHelloName, _ = bareHelloName.AppendComponent(randomSuffix)

	// Send the bare hello
	log.Println("Sending the hello")
	bareHello := kex.KEXHello()
	bareHelloRequest := interest.CreateWithName(bareHelloName)
	bareHelloRequest.AddContainer(bareHello)
	n.apiStack.Push(messages.Package(bareHelloRequest))

	// Wait for the response, and use it to build the full hello
	replyWrapper := n.apiStack.Pop()
	reply := replyWrapper.InnerMessage()
	log.Println("Got the REJECT")

	reject, err := reply.GetContainer(codec.T_KEX)
	if err != nil {
		log.Println("Error: no KEX container in the REJECT content object")
		return
	}
	hello := kex.KEXFullHello(bareHello, reject.(*kex.KEX))

	randomSuffix, _ = random.GenerateRandomString(16)
	helloName, _ := prefix.AppendComponent(connectString)
	helloName, _ = prefix.AppendComponent(randomSuffix)

	helloRequest := interest.CreateWithName(helloName)
	helloRequest.AddContainer(hello)
	n.apiStack.Push(messages.Package(helloRequest))

	// Wait for the response to complete the KEX
	replyWrapper = n.apiStack.Pop()
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

	session := tunnel.NewSession(sharedKey[:], acceptKEX.GetSessionID())
	n.apiStack.AddSession(session, prefix)

	log.Println("Consumer key: ", sharedKey)
}

func (n SecurePortal) Get(request *messages.MessageWrapper, timeout time.Duration) (*messages.MessageWrapper, error) {
	signalChannel := make(chan *messages.MessageWrapper, 1)
	n.apiStack.Get(request, func(msg *messages.MessageWrapper) {
		signalChannel <- msg
	})

	var response *messages.MessageWrapper
	select {
	case data := <-signalChannel:
		return data, nil
	case <-time.After(timeout):
		return response, PortalError{0, "Timeout"}
	}
}

func (n SecurePortal) GetAsync(request *messages.MessageWrapper, callback ResponseMessageCallback) {
	n.apiStack.Get(request, func(msg *messages.MessageWrapper) {
		callback(msg)
	})
}

func (p SecurePortal) GetAsyncWithTimeout(request *messages.MessageWrapper, timeout time.Duration, callback ResponseMessageCallback) {
	signalChannel := make(chan *messages.MessageWrapper, 1)
	p.apiStack.Get(request, func(msg *messages.MessageWrapper) {
		signalChannel <- msg
	})

	select {
	case data := <-signalChannel:
		callback(data)
	case <-time.After(timeout):
		p.apiStack.Cancel(request)
	}
}

func (n SecurePortal) Serve(prefix *name.Name, callback RequestMessageCallback) {
	if prefix == nil {
		return
	}

	established := false
	for {
		requestWrapper := n.apiStack.Pop()

		if established {
			log.Println("Handling a request")
			response := callback(requestWrapper)
			if response != nil {
				n.apiStack.Push(response)
			} else {
				log.Println("Failed to generate a response")
			}
		} else {
			request := requestWrapper.InnerMessage()
			if !prefix.IsPrefix(request.Name()) {
				break
			}

			kexTLV, _ := request.GetContainer(codec.T_KEX)
			kexContainer := kexTLV.(*kex.KEX)

			switch kexContainer.GetMessageType() {
			case codec.T_KEX_BAREHELLO:
				log.Println("Got the BARE HELLO")
				reject := kex.KEXHelloReject(kexContainer, n.macKey)
				rejectResponse := content.CreateWithName(request.Name())
				rejectResponse.AddContainer(reject)
				n.apiStack.Push(messages.Package(rejectResponse))
				break

			case codec.T_KEX_HELLO:
				log.Println("Got the HELLO")
				accept, err := kex.KEXHelloAccept(kexContainer, n.macKey, n.encKey)
				if err != nil {
					log.Println(err)
					break
				}

				acceptResponse := content.CreateWithName(request.Name())
				acceptResponse.AddContainer(accept)
				n.apiStack.Push(messages.Package(acceptResponse))

				// XXX: go to the KDF step

				var sharedKey [32]byte
				var peerPublic [32]byte
				var privateKey [32]byte
				copy(peerPublic[:], kexContainer.GetPublicKeyShare())
				copy(privateKey[:], accept.GetPrivateKeyShare())
				box.Precompute(&sharedKey, &peerPublic, &privateKey)

				log.Println("Producer key:", sharedKey)

				// Create and start the session
				// session := esic.NewESIC(n.apiStack, sharedKey[:], accept.GetSessionID())
				// callback(session)
				time.Sleep(100 * time.Millisecond)

				session := tunnel.NewSession(sharedKey[:], accept.GetSessionID())
				n.apiStack.AddSession(session, prefix)
				established = true

				break

			case codec.T_KEX_REJECT:
			case codec.T_KEX_ACCEPT:
				log.Println("Got an invalid message...")
				// invalid message type to be received here...
				break
			}
		}
	}
}

func (p SecurePortal) Produce(data *messages.MessageWrapper) {
	p.apiStack.Push(data)
}
