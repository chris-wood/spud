package crypto

import "fmt"
import "log"

import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/interest"
import "github.com/chris-wood/spud/messages/validation"
import "github.com/chris-wood/spud/messages/validation/publickey"
import "github.com/chris-wood/spud/stack/component"
import "github.com/chris-wood/spud/stack/component/crypto/validator"
import "github.com/chris-wood/spud/stack/component/crypto/context"

type CryptoComponent struct {
	ingress chan *messages.MessageWrapper
	egress  chan *messages.MessageWrapper

	// Pending egress messages
	pendingMap map[string]*messages.MessageWrapper

	// XXX: queue of packets pending verification
	pendingVerificationQueue map[string]*messages.MessageWrapper

	// The context will be modified
	trustStore *context.TrustStore
	downstream component.Component

	// XXX: LPM table of processors
	cryptoProcessor validator.CryptoProcessor
}

func NewCryptoComponent(trustStore *context.TrustStore, downstream component.Component) CryptoComponent {
	egress := make(chan *messages.MessageWrapper)
	ingress := make(chan *messages.MessageWrapper)

	return CryptoComponent{
		ingress:                  ingress,
		egress:                   egress,
		trustStore:               trustStore,
		downstream:               downstream,
		pendingMap:               make(map[string]*messages.MessageWrapper),
		pendingVerificationQueue: make(map[string]*messages.MessageWrapper),
	}
}

func (c *CryptoComponent) AddCryptoProcessor(pattern string, proc validator.CryptoProcessor) {
	c.cryptoProcessor = proc
	validationAlgorithm := proc.ProcessorDetails()

	// XXX: build the named key (certificate) and add it to the cache

	c.trustStore.AddTrustedKey(validationAlgorithm.KeyIdString(), validationAlgorithm.GetPublicKey())
}

func addAuthenticator(msg *messages.MessageWrapper, proc validator.CryptoProcessor) (*messages.MessageWrapper, error) {
	va := proc.ProcessorDetails()
	msg.SetValidationAlgorithm(va)
	signature, err := proc.Sign(msg)
	if err == nil {
		vp := validation.NewValidationPayload(signature)
		msg.SetValidationPayload(vp)
	} else {
		log.Println("Signing error:", err)
	}

	return msg, err
}

func (c CryptoComponent) ProcessEgressMessages() {
	for {
		msg := <-c.egress

		// Look up the processor based on the message
		// XXX: apply the LPM filter for the right processor here
		// XXX: processor is identified by the name only
		var err error
		msg, err = addAuthenticator(msg, c.cryptoProcessor)
		if err != nil {
			fmt.Println(err.Error())
		}
		c.downstream.Enqueue(msg)

		// Save a reocrd for all egress requests
		if msg.GetPacketType() == codec.T_INTEREST {
			c.pendingMap[msg.Identifier()] = msg
		}
	}
}

func (c CryptoComponent) handleIngressRequest(msg *messages.MessageWrapper) {
	c.ingress <- msg
}

// Check to see if there are any other messages to verify with this
// newly verified key. If there is, recursively call the verify request
func (c *CryptoComponent) processPendingResponses(msg *messages.MessageWrapper) {
	dependentRequest, ok := c.pendingVerificationQueue[msg.Identifier()]
	if ok {
		c.handleIngressResponse(dependentRequest)
		delete(c.pendingVerificationQueue, msg.Identifier())

		if msg.PayloadType() == codec.T_PAYLOADTYPE_KEY {
			payload := msg.Payload()
			rawKey := publickey.Create(payload.Value())

			c.trustStore.AddTrustedKey(rawKey.KeyIdString(), rawKey)
		}
	} else {
		c.ingress <- msg
		delete(c.pendingMap, msg.Identifier())
	}
}

func (c CryptoComponent) dropPendingResponses(msg *messages.MessageWrapper) {
	delete(c.pendingVerificationQueue, msg.Identifier())
	delete(c.pendingMap, msg.Identifier())
}

// XXX: rename this
func (c CryptoComponent) checkTrustProperties(msg *messages.MessageWrapper) {
	validationAlgorithm := msg.GetValidationAlgorithm()
	keyId := validationAlgorithm.KeyIdString()

	if c.trustStore.IsTrustedKey(keyId) {
		c.processPendingResponses(msg)
	} else {
		log.Println("Not a trusted key. Drop the message.")
		c.dropPendingResponses(msg)
	}
}

func (c CryptoComponent) handleIngressResponse(msg *messages.MessageWrapper) {
	// Check to see if this is a response to a previous key name
	request, isPending := c.pendingMap[msg.Identifier()]
	if isPending {

		// XXX: how to identify the right processor? based on the name only?

		if !c.cryptoProcessor.CanVerify(msg) {
			// Pull out the key name
			// XXX: here we'd swtich on the type of locator
			va := msg.GetValidationAlgorithm()
			link := va.GetKeyLink()

			// Create an interest for the link and send it
			keyMsg := interest.CreateFromLink(link)
			keyPacket := messages.Package(keyMsg)

			c.downstream.Enqueue(keyPacket)
			c.pendingMap[keyPacket.Identifier()] = keyPacket

			// Save the reference to this response
			c.pendingVerificationQueue[keyPacket.Identifier()] = msg
		} else {
			if c.cryptoProcessor.Verify(request, msg) {
				c.checkTrustProperties(msg)
			} else {
				log.Println("Dropping the response")
				c.dropPendingResponses(msg)
			}
		}
	} else {
		log.Println("Error: no matching request found: " + msg.Identifier())
		log.Println("Dropping the packet.")
	}
}

func (c CryptoComponent) ProcessIngressMessages() {
	for {
		msg := c.downstream.Dequeue()

		// Hand off the message to the request/response handler
		if msg.GetPacketType() != codec.T_INTEREST {
			go c.handleIngressResponse(msg)
		} else {
			go c.handleIngressRequest(msg)
		}
	}
}

func (c CryptoComponent) Enqueue(msg *messages.MessageWrapper) {
	c.egress <- msg
}

func (c CryptoComponent) Dequeue() *messages.MessageWrapper {
	msg := <-c.ingress
	return msg
}
