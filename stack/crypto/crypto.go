package crypto

import "fmt"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/stack/codec"
import "github.com/chris-wood/spud/stack/crypto/processor"

type CryptoComponent struct {
    ingress chan messages.Message
    egress chan messages.Message

    cryptoProcessor processor.CryptoProcessor
    codecComponent codec.CodecComponent
}

func NewCryptoComponent(cryptoProcessor processor.CryptoProcessor, codecComponent codec.CodecComponent) CryptoComponent {
    egress := make(chan messages.Message)
    ingress := make(chan messages.Message)

    return CryptoComponent{ingress: ingress, egress: egress, cryptoProcessor: cryptoProcessor, codecComponent: codecComponent}
}

func (c CryptoComponent) ProcessEgressMessages() {
    for ;; {
        msg := <- c.egress
        fmt.Println("Passing down: " + msg.Identifier())

        // 0. Look up the processor based on the message, and then extract its validation algorithm

        // 1. Add the key locator information
        // XXX

        // 2. Hash the sensitive region
        // XXX

        // 3. Compute the signature
        // XXX

        // 4. Append the signature
        // XXX

        c.codecComponent.Enqueue(msg)
    }
}

func (c CryptoComponent) ProcessIngressMessages() {
    for ;; {
        msg := c.codecComponent.Dequeue()

        fmt.Println("Passing up: " + msg.Identifier())

        // 1. Hash the sensitive region
        // XXX

        // 2. Verify the signature
        // XXX

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
