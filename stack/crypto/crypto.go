package crypto

import "fmt"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/validation"
import "github.com/chris-wood/spud/stack/codec"
import "github.com/chris-wood/spud/stack/crypto/processor"

type CryptoComponent struct {
    ingress chan messages.Message
    egress chan messages.Message

    pendingMap map[string]messages.Message

    cryptoProcessor processor.CryptoProcessor
    codecComponent codec.CodecComponent
}

func NewCryptoComponent(cryptoProcessor processor.CryptoProcessor, codecComponent codec.CodecComponent) CryptoComponent {
    egress := make(chan messages.Message)
    ingress := make(chan messages.Message)

    return CryptoComponent{ingress: ingress, egress: egress, cryptoProcessor: cryptoProcessor, codecComponent: codecComponent, pendingMap: make(map[string]messages.Message)}
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

// XXX: this does not look at the request, yet
// func verifyAuthenticator(request, response messages.Message, crypto processor.CryptoProcessor) (bool, error) {
//     validationPayload := response.GetValidationPayload()
//     validationAlgorithm := response.GetValidationAlgorithm()
//
//     switch validationAlgorithm.Type() {
//     case codec.T_RSA_SHA256:
//     default:
//         return false,
//     }
//
//     signature := validationPayload.Value()
//     valid := crypto.Verify(response, signature)
//     return valid, nil
// }

func (c CryptoComponent) ProcessEgressMessages() {
    for ;; {
        msg := <- c.egress
        fmt.Println("Passing down: " + msg.Identifier())

        // Look up the processor based on the message
        // XXX: apply the LPM filter for the right processor here
        // if !msg.IsRequest() {
            msg, err := addAuthenticator(msg, c.cryptoProcessor)
            if err != nil {
                fmt.Println(err.Error())
            }
            c.codecComponent.Enqueue(msg)
        // }

        if msg.IsRequest() {
            c.pendingMap[msg.Identifier()] = msg
            // c.codecComponent.Enqueue(msg)
        }
    }
}

func (c CryptoComponent) ProcessIngressMessages() {
    for ;; {
        msg := c.codecComponent.Dequeue()
        fmt.Println("Passing up: " + msg.Identifier())

        // XXX: clean up this logic
        if !msg.IsRequest() {
            request, ok := c.pendingMap[msg.Identifier()]
            if ok {
                // XXX: this should return an error
                success := c.cryptoProcessor.Verify(request, msg)
                if success {
                    fmt.Println(">>> valid signature.")
                } else {
                    fmt.Println("invalid signature...")
                }
            } else {
                fmt.Println("Drop the message: " + msg.Identifier())
            }
        }

        // 3. If valid, enqueue upstream
        c.ingress <- msg

        // 4. Else, request whatever is needed to verify the signature and keep going
    }
}

func (c CryptoComponent) Enqueue(msg messages.Message) {
    c.egress <- msg
}

func (c CryptoComponent) Dequeue() messages.Message {
    msg := <-c.ingress
    return msg
}
