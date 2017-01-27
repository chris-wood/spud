package tunnel

import "log"

import messageCodec "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/messages/interest"
import "github.com/chris-wood/spud/messages/content"
import "github.com/chris-wood/spud/stack/component"
import "github.com/chris-wood/spud/stack/component/codec"

type Tunnel struct {
	ingress chan messages.MessageWrapper
	egress  chan messages.MessageWrapper
    downstream component.Component
    session *Session
}

type TunnelComponent struct {
	ingress chan messages.MessageWrapper
	egress  chan messages.MessageWrapper
    exitCodec *codec.CodecComponent
	tunnels []*Tunnel
}

func NewTunnel(session *Session, downstream component.Component) *Tunnel {
	egress := make(chan messages.MessageWrapper)
	ingress := make(chan messages.MessageWrapper)

	log.Println("Created Transport component")

	return &Tunnel{ingress: ingress, egress: egress, downstream: downstream, session: session}
}

func (c Tunnel) ProcessEgressMessages() {
	for {
		msg := <-c.egress

		encodedRequest := msg.Encode()
		encryptedMessage, err := c.session.Encrypt(encodedRequest)
		if err == nil {
			baseName := msg.Name()
			sessionName, _ := baseName.AppendComponent(c.session.SessionID)

			encapPayload := payload.Create(encryptedMessage)

			var encapResponse messages.MessageWrapper
			if msg.GetPacketType() == messageCodec.T_INTEREST {
				encapResponse = messages.Package(interest.CreateWithNameAndPayload(sessionName, messageCodec.T_PAYLOADTYPE_ENCAP, encapPayload))
			} else {
				encapResponse = messages.Package(content.CreateWithNameAndTypedPayload(sessionName, messageCodec.T_PAYLOADTYPE_ENCAP, encapPayload))
			}

			c.egress <- encapResponse
			// c.codecComponent.Enqueue(encapResponse)
		}
	}
}

func (c Tunnel) ProcessIngressMessages() {
	for {
		// msg := c.codecComponent.Dequeue()
		var msg messages.MessageWrapper
		encryptedPayload := msg.InnerMessage().Payload().Value()
		encapInterest, err := c.session.Decrypt(encryptedPayload)
		if err == nil {
			d := messageCodec.Decoder{}
			decodedTlV := d.Decode(encapInterest)
			if len(decodedTlV) == 1 {
				responseMsg, err := messages.CreateFromTLV(decodedTlV)
				if err == nil {
					c.ingress <- responseMsg
				}
			}
		}
	}
}

func (c Tunnel) Enqueue(msg messages.MessageWrapper) {
    c.egress <- msg
}

func (c Tunnel) Dequeue() messages.MessageWrapper {
    msg := <-c.ingress
    return msg
}

func NewTunnelComponent(codecComponent *codec.CodecComponent) TunnelComponent {
	egress := make(chan messages.MessageWrapper)
	ingress := make(chan messages.MessageWrapper)

	log.Println("Created Transport component")

	return TunnelComponent{ingress: ingress, egress: egress, tunnels: make([]*Tunnel, 0)}
}

func (c TunnelComponent) ProcessEgressMessages() {
	for {
		msg := <-c.egress
        if len(c.tunnels) == 0 {
            c.exitCodec.Enqueue(msg)
        } else {
            c.tunnels[0].Enqueue(msg)
        }
	}
}

func (c TunnelComponent) ProcessIngressMessages() {
	for {
        if len(c.tunnels) == 0 {
            msg := c.exitCodec.Dequeue()
            c.ingress <- msg
        } else {
            msg := c.tunnels[0].Dequeue()
            c.ingress <- msg
        }
	}
}

func (c TunnelComponent) Enqueue(msg messages.MessageWrapper) {
    if len(c.tunnels) == 0 {
        c.egress <- msg
    } else {
        c.tunnels[0].Enqueue(msg)
    }
}

func (c TunnelComponent) Dequeue() messages.MessageWrapper {
    msg := <-c.ingress
    return msg
}
