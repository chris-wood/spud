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
	"fmt"
)

type Tunnel struct {
	ingress    chan *messages.MessageWrapper
	egress     chan *messages.MessageWrapper
	baseName   *name.Name
	downstream component.Component
	session    *Session

	// XXX: still needed?
	nameMap map[*name.Name]*name.Name
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
	t := Tunnel{ingress: ingress, egress: egress, baseName: baseName, downstream: downstream, session: session, nameMap: make(map[*name.Name]*name.Name)}
	return &t
}

func (c *Tunnel) ProcessEgressMessages() {
	for {
		msg := <-c.egress

		encodedRequest := msg.Encode()

		log.Println("Encrypting", msg.Name())
		//fmt.Println(encodedRequest)

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

			c.nameMap[sessionName] = msg.Name()
			c.downstream.Push(encapResponse)
		}
	}
}

func (c *Tunnel) ProcessIngressMessages() {
	for {
		msg := c.downstream.Pop()
		encryptedPayload := msg.InnerMessage().Payload().Value()
		encapInterest, err := c.session.Decrypt(encryptedPayload)

		log.Println("Decrypting", msg.Identifier())

		//fmt.Println("DECAPSULATED INTEREST", encapInterest)

		if err == nil {
			d := codec.Decoder{}
			decodedTlV := d.Decode(encapInterest)
			responseMsg, err := messages.CreateFromTLV(decodedTlV)

			//log.Println(responseMsg.Identifier())
			//log.Println(responseMsg.InnerMessage().Payload())
			//log.Println(msg.Identifier())

			if err == nil {
				fmt.Println("Sending up decoded,", responseMsg)
				c.ingress <- responseMsg
			}
		} else {
			log.Println("Failed to decrypt packet")
		}
	}
}

func (c *Tunnel) Inject(msg *messages.MessageWrapper) {
	c.ingress <- msg
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
		go tunnel.ProcessEgressMessages()
		go tunnel.ProcessIngressMessages()
		c.tunnels = []*Tunnel{tunnel}
	} else {
		tunnel := NewTunnel(session, baseName, c.tunnels[0])
		go tunnel.ProcessEgressMessages()
		go tunnel.ProcessIngressMessages()
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
		//buffer := make(chan *messages.MessageWrapper)

		if len(c.tunnels) == 0 {
			msg := c.exitCodec.Pop()
			if len(c.tunnels) > 0 {
				fmt.Println("Reinstering into the tunnel")
				c.tunnels[len(c.tunnels) - 1].downstream.Inject(msg)
			} else {
				c.ingress <- msg
			}
		} else {
			msg := c.tunnels[0].Pop()
			c.ingress <- msg
		}

		// XXX: the exit codec is stealing the message from the tunnel... causing it to not be decrypted, and things to fail
		//go func () {
		//	buffer <- c.exitCodec.Pop()
		//	log.Println("Received from codec")
		//}()
		//go func () {
		//	if len(c.tunnels) > 0 {
		//		buffer <- c.tunnels[0].Pop()
		//		log.Println("Received from tunnel")
		//	}
		//}()
		//
		//msg := <-buffer
		//c.ingress <- msg

		//select {
		//case plainMsg := c.exitCodec.Pop():
		//	log.Println(len(c.tunnels))
		//	log.Println("Processing ingress", plainMsg.Name())
		//	c.ingress <- plainMsg
		//case tunneledMsg := c.tunnels[0].Pop():
		//	log.Println(len(c.tunnels))
		//	log.Println("Processing tunneled ingress", tunneledMsg.Name())
		//	c.ingress <- tunneledMsg
		//}
	}
}

func (c *TunnelComponent) Inject(msg *messages.MessageWrapper) {
	c.ingress <- msg
}

func (c *TunnelComponent) Push(msg *messages.MessageWrapper) {
	c.egress <- msg
}

func (c *TunnelComponent) Pop() *messages.MessageWrapper {
	msg := <-c.ingress
	return msg
}
