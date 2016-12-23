package crypto

import "fmt"
import tlvCodec "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/interest"
import "github.com/chris-wood/spud/messages/validation"
import "github.com/chris-wood/spud/messages/validation/publickey"
import "github.com/chris-wood/spud/stack/component/codec"
import "github.com/chris-wood/spud/stack/component/crypto/processor"
import "github.com/chris-wood/spud/stack/component/crypto/context"

type CryptoComponent struct {
    ingress chan messages.Message
    egress chan messages.Message

    // Pending egress messages
    pendingMap map[string]messages.Message

    // XXX: queue of packets pending verification
    pendingVerificationQueue map[string]messages.Message

    // The context will be modified
    cryptoContext *context.CryptoContext
    codecComponent codec.CodecComponent

    // XXX: LPM table of processors
    cryptoProcessor processor.CryptoProcessor
}

func NewCryptoComponent(cryptoContext *context.CryptoContext, codecComponent codec.CodecComponent) CryptoComponent {
    egress := make(chan messages.Message)
    ingress := make(chan messages.Message)

    return CryptoComponent{
        ingress: ingress,
        egress: egress,
        cryptoContext: cryptoContext,
        codecComponent: codecComponent,
        pendingMap: make(map[string]messages.Message),
        pendingVerificationQueue: make(map[string]messages.Message),
    }
}

func (c *CryptoComponent) AddCryptoProcessor(pattern string, proc processor.CryptoProcessor) {
    c.cryptoProcessor = proc
    validationAlgorithm := proc.ProcessorDetails()

    // XXX: build the named key (certificate) and add it to the cache

    c.cryptoContext.AddTrustedKey(validationAlgorithm.KeyIdString(), validationAlgorithm.GetPublicKey())
}

func addAuthenticator(msg messages.Message, proc processor.CryptoProcessor) (messages.Message, error) {
    va := proc.ProcessorDetails()
    msg.SetValidationAlgorithm(va)
    signature, err := proc.Sign(msg)
    if err != nil {
        vp := validation.NewValidationPayload(signature)
        msg.SetValidationPayload(vp)
    }

    return msg, err
}

func (c CryptoComponent) ProcessEgressMessages() {
    for ;; {
        msg := <- c.egress

        // Look up the processor based on the message
        // XXX: apply the LPM filter for the right processor here
        // XXX: processor is identified by the name only
        // if !msg.IsRequest() {

        var err error
        msg, err = addAuthenticator(msg, c.cryptoProcessor)
        if err != nil {
            fmt.Println(err.Error())
        }
        c.codecComponent.Enqueue(msg)
        // }

        // XXX: move this code to a function
        if msg.GetPacketType() == tlvCodec.T_INTEREST {
            c.pendingMap[msg.Identifier()] = msg
            // c.codecComponent.Enqueue(msg)
        }
    }
}

func (c CryptoComponent) handleIngressRequest(msg messages.Message) {
    // XXX: what else needs to be done here?
    c.ingress <- msg
}

// Check to see if there are any other messages to verify with this
// newly verified key. If there is, recursively call the verify request
func (c *CryptoComponent) processPendingResponses(msg messages.Message) {
    dependentRequest, ok := c.pendingVerificationQueue[msg.Identifier()]
    if ok {
        c.handleIngressResponse(dependentRequest)
        delete(c.pendingVerificationQueue, msg.Identifier())

        if msg.PayloadType() == tlvCodec.T_PAYLOADTYPE_KEY {
            payload := msg.Payload()
            rawKey := publickey.Create(payload.Value())

            c.cryptoContext.AddTrustedKey(rawKey.KeyIdString(), rawKey)
        }
    } else {
        c.ingress <- msg
//        fmt.Println("Dropping pending response:", msg.Identifier())
//        fmt.Println(c.pendingMap)
//        fmt.Println(msg.GetPacketType() == tlvCodec.T_INTEREST)
        delete(c.pendingMap, msg.Identifier())
    }
}

func (c CryptoComponent) dropPendingResponses(msg messages.Message) {
    delete(c.pendingVerificationQueue, msg.Identifier())
    delete(c.pendingMap, msg.Identifier())
}

// XXX: rename this
func (c CryptoComponent) checkTrustProperties(msg messages.Message) {
    validationAlgorithm := msg.GetValidationAlgorithm()
    keyId := validationAlgorithm.KeyIdString()

    if c.cryptoContext.IsTrustedKey(keyId) {
        c.processPendingResponses(msg)
    } else {
        fmt.Println("Not a trusted key. Abort.")
        c.dropPendingResponses(msg)
    }
}

func (c CryptoComponent) handleIngressResponse(msg messages.Message) {
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
            c.codecComponent.Enqueue(keyMsg)
            c.pendingMap[keyMsg.Identifier()] = keyMsg

            // Save the reference to this response
            c.pendingVerificationQueue[keyMsg.Identifier()] = msg
        } else {
            if c.cryptoProcessor.Verify(request, msg) {
                c.checkTrustProperties(msg)
            } else {
                c.dropPendingResponses(msg)
            }
        }
    } else {
        fmt.Println("Error: no matching request found: " + msg.Identifier())
        fmt.Println("Dropping the packet.")
    }
}

func (c CryptoComponent) ProcessIngressMessages() {
    for ;; {
        msg := c.codecComponent.Dequeue()

        // Hand off the message to the request/response handler
        if msg.GetPacketType() != tlvCodec.T_INTEREST {
            go c.handleIngressResponse(msg)
        } else {
            go c.handleIngressRequest(msg)
        }
    }
}

func (c CryptoComponent) Enqueue(msg messages.Message) {
    c.egress <- msg
}

func (c CryptoComponent) Dequeue() messages.Message {
    msg := <-c.ingress
    return msg
}
