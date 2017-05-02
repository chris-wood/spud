package tunnel

import (
	"github.com/chris-wood/spud/codec"
	"github.com/chris-wood/spud/messages"
	"github.com/chris-wood/spud/messages/name"
	"github.com/chris-wood/spud/stack/component"
	"github.com/chris-wood/spud/messages/payload"
	"github.com/chris-wood/spud/messages/interest"
	"github.com/chris-wood/spud/messages/content"
	"log"
)

type Tunnel struct {
	ingress    chan *messages.MessageWrapper
	egress     chan *messages.MessageWrapper
	baseName   *name.Name
	downstream component.Component
	session    *Session
}

type TunnelComponent struct {
	tunnelNames map[*name.Name]bool
	ingress   chan *messages.MessageWrapper
	egress    chan *messages.MessageWrapper
	exitCodec component.Component
	tunnels   []*Tunnel
}

func NewTunnel(session *Session, baseName *name.Name, downstream component.Component) *Tunnel {
	egress := make(chan *messages.MessageWrapper)
	ingress := make(chan *messages.MessageWrapper)
	t := Tunnel{ingress: ingress, egress: egress, baseName: baseName, downstream: downstream, session: session}
	return &t
}

func (c *Tunnel) ProcessEgressMessages() {
	for {
		msg := <-c.egress

		encodedRequest := msg.Encode()
		encryptedMessage, err := c.session.Encrypt(encodedRequest)
		if err == nil {
			sessionName, _ := c.baseName.AppendComponent(c.session.SessionID)
			// TODO(cawood): append the packet counter, too

			encapPayload := payload.Create(encryptedMessage)

			var encapResponse *messages.MessageWrapper
			if msg.GetPacketType() == codec.T_INTEREST {
				encapResponse = messages.Package(interest.CreateWithNameAndPayload(sessionName, codec.T_PAYLOADTYPE_ENCAP, encapPayload))
			} else {
				encapResponse = messages.Package(content.CreateWithNameAndTypedPayload(sessionName, codec.T_PAYLOADTYPE_ENCAP, encapPayload))
			}

			c.downstream.Push(encapResponse)
		}
	}
}

func (c *Tunnel) ProcessIngressMessages() {
	for {
		msg := c.downstream.Pop()
		encryptedPayload := msg.InnerMessage().Payload().Value()
		encapInterest, err := c.session.Decrypt(encryptedPayload)
		if err == nil {
			d := codec.Decoder{}
			decodedTlV := d.Decode(encapInterest)
			responseMsg, err := messages.CreateFromTLV(decodedTlV)

			//responseMsg.Name().DropSuffix() // Drop the session ID from the name

			if err == nil {
				c.ingress <- responseMsg
			}
		} else {
			log.Println("Failed to decrypt packet")
		}
	}
}

func (c *Tunnel) Push(msg *messages.MessageWrapper) {
	c.egress <- msg
}

func (c *Tunnel) Pop() *messages.MessageWrapper {
	msg := <-c.ingress
	return msg
}

func NewTunnelComponent(exitComponent component.Component) *TunnelComponent {
	egress := make(chan *messages.MessageWrapper)
	ingress := make(chan *messages.MessageWrapper)
	return &TunnelComponent{tunnelNames: make(map[*name.Name]bool), ingress: ingress, egress: egress, exitCodec: exitComponent, tunnels: make([]*Tunnel, 0)}
}

func (c *TunnelComponent) AddSession(session *Session, baseName *name.Name) {
	if _, ok := c.tunnelNames[baseName]; ok {
		log.Println("Tunnel already exists.")
		return
	}


	if len(c.tunnels) == 0 {
		tunnel := NewTunnel(session, baseName, c.exitCodec)
		go tunnel.ProcessIngressMessages()
		go tunnel.ProcessEgressMessages()
		c.tunnels = []*Tunnel{tunnel}
	} else {
		tunnel := NewTunnel(session, baseName, c.tunnels[0])
		go tunnel.ProcessIngressMessages()
		go tunnel.ProcessEgressMessages()
		c.tunnels = append([]*Tunnel{tunnel}, c.tunnels...)
	}

	// Add the tunnel to the list
	c.tunnelNames[baseName] = true
}

func (c *TunnelComponent) ProcessEgressMessages() {
	for {
		msg := <-c.egress
		if len(c.tunnels) == 0 {
			c.exitCodec.Push(msg)
		} else {
			c.tunnels[0].Push(msg)
		}
	}
}

func (c *TunnelComponent) ProcessIngressMessages() {
	for {
		if len(c.tunnels) == 0 {
			msg := c.exitCodec.Pop()
			c.ingress <- msg
		} else {
			msg := c.tunnels[0].Pop()
			c.ingress <- msg
		}
	}
}

func (c *TunnelComponent) Push(msg *messages.MessageWrapper) {
	c.egress <- msg
}

func (c *TunnelComponent) Pop() *messages.MessageWrapper {
	msg := <-c.ingress
	return msg
}
