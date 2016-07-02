package crypto

import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/stack/codec"

type CryptoComponent struct {
    ingress chan messages.Message
    egress chan messages.Message

    codecComponent codec.CodecComponent
}

func NewCryptoComponent(codecComponent codec.CodecComponent) CryptoComponent {
    egress := make(chan messages.Message)
    ingress := make(chan messages.Message)

    return CryptoComponent{ingress: ingress, egress: egress, codecComponent: codecComponent}
}

type CryptoProcessor interface {
    Sign(msg *messages.Message) []byte
    Verify(msg *messages.Message) bool
}

type XXXProcessor struct {
    // XX key store
}

func (c CryptoComponent) ProcessEgressMessages() {
    for ;; {
        msg := <- c.egress

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

        // 1. Hash the sensitive region
        // XXX

        // 2. Verify the signature
        // XXX

        // 3. If valid, enqueue upstream
        c.ingress <- message

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
