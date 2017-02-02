package tunnel

import "log"

import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/name"
import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/messages/interest"
import "github.com/chris-wood/spud/messages/content"
import "github.com/chris-wood/spud/stack/component"

type Tunnel struct {
	ingress chan *messages.MessageWrapper
	egress  chan *messages.MessageWrapper
    baseName *name.Name
    downstream component.Component
    session *Session
}

type TunnelComponent struct {
	ingress chan *messages.MessageWrapper
	egress  chan *messages.MessageWrapper
    exitCodec component.Component
	tunnels []*Tunnel
}

func NewTunnel(session *Session, baseName *name.Name, downstream component.Component) *Tunnel {
	egress := make(chan *messages.MessageWrapper)
	ingress := make(chan *messages.MessageWrapper)

	log.Println("Created Tunnel")

	t := Tunnel{ingress: ingress, egress: egress, baseName: baseName, downstream: downstream, session: session}
    return &t
}

func (c Tunnel) ProcessEgressMessages() {
	for {
		msg := <-c.egress

		encodedRequest := msg.Encode()
		encryptedMessage, err := c.session.Encrypt(encodedRequest)
		if err == nil {
			sessionName, _ := c.baseName.AppendComponent(c.session.SessionID)
            // XXX: append the packet counter, too

			encapPayload := payload.Create(encryptedMessage)

			var encapResponse *messages.MessageWrapper
			if msg.GetPacketType() == codec.T_INTEREST {
				encapResponse = messages.Package(interest.CreateWithNameAndPayload(sessionName, codec.T_PAYLOADTYPE_ENCAP, encapPayload))
			} else {
				encapResponse = messages.Package(content.CreateWithNameAndTypedPayload(sessionName, codec.T_PAYLOADTYPE_ENCAP, encapPayload))
			}

			c.downstream.Enqueue(encapResponse)
		}
	}
}

func (c Tunnel) ProcessIngressMessages() {
	for {
		msg := c.downstream.Dequeue()
		encryptedPayload := msg.InnerMessage().Payload().Value()
		encapInterest, err := c.session.Decrypt(encryptedPayload)
		if err == nil {
			d := codec.Decoder{}
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

func (c Tunnel) Enqueue(msg *messages.MessageWrapper) {
    c.egress <- msg
}

func (c Tunnel) Dequeue() *messages.MessageWrapper {
    msg := <-c.ingress
    return msg
}

func NewTunnelComponent(exitComponent component.Component) *TunnelComponent {
	egress := make(chan *messages.MessageWrapper)
	ingress := make(chan *messages.MessageWrapper)

	log.Println("Created TunnelComponent")

	return &TunnelComponent{ingress: ingress, egress: egress, exitCodec: exitComponent, tunnels: make([]*Tunnel, 0)}
}

func (c *TunnelComponent) AddSession(session *Session, baseName *name.Name) {
    if len(c.tunnels) == 0 {
        tunnel := NewTunnel(session, baseName, c.exitCodec)
        go tunnel.ProcessEgressMessages()
        go tunnel.ProcessIngressMessages()
        c.tunnels = append(c.tunnels, tunnel)
    } else {
        tunnel := NewTunnel(session, baseName, c.tunnels[0])
        go tunnel.ProcessEgressMessages()
        go tunnel.ProcessIngressMessages()
        c.tunnels = append([]*Tunnel{tunnel}, c.tunnels...)
    }
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

func (c TunnelComponent) Enqueue(msg *messages.MessageWrapper) {
    c.egress <- msg
}

func (c TunnelComponent) Dequeue() *messages.MessageWrapper {
    msg := <-c.ingress
    return msg
}
