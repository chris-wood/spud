package transport

import "log"

import "github.com/chris-wood/spud/messages"

type TransportComponent struct {
    ingress chan messages.MessageWrapper
    egress chan messages.MessageWrapper
}

func NewTransportComponent() TransportComponent {
    egress := make(chan messages.MessageWrapper)
    ingress := make(chan messages.MessageWrapper)

    log.Println("Created Transport component")

    return TransportComponent{ingress: ingress, egress: egress}
}

func (c TransportComponent) ProcessEgressMessages() {
    for ;; {
        // XXX
    }
}

func (c TransportComponent) ProcessIngressMessages() {
    for ;; {
        // XXX
    }
}

func (c TransportComponent) Enqueue(msg messages.MessageWrapper) {
    c.egress <- msg
}

func (c TransportComponent) Dequeue() messages.MessageWrapper {
    msg := <-c.ingress
    return msg
}
