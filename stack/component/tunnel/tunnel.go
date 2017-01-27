package tunnel

import "log"

import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/messages/payload"
import "github.com/chris-wood/spud/messages/interest"
import "github.com/chris-wood/spud/messages/content"

type TunnelComponent struct {
	ingress chan messages.MessageWrapper
	egress  chan messages.MessageWrapper
	session *Session
}

func NewTunnelComponent(session *Session) TunnelComponent {
	egress := make(chan messages.MessageWrapper)
	ingress := make(chan messages.MessageWrapper)

	log.Println("Created Transport component")

	return TunnelComponent{ingress: ingress, egress: egress, session: session}
}

func (c TunnelComponent) ProcessEgressMessages() {
	for {
		msg := <-c.egress

		encodedRequest := msg.Encode()
		encryptedMessage, err := c.session.Encrypt(encodedRequest)
		if err == nil {
			baseName := msg.Name()
			sessionName, _ := baseName.AppendComponent(c.session.SessionID)

			encapPayload := payload.Create(encryptedMessage)

			var encapResponse messages.MessageWrapper
			if msg.GetPacketType() == codec.T_INTEREST {
				encapResponse = messages.Package(interest.CreateWithNameAndPayload(sessionName, codec.T_PAYLOADTYPE_ENCAP, encapPayload))
			} else {
				encapResponse = messages.Package(content.CreateWithNameAndTypedPayload(sessionName, codec.T_PAYLOADTYPE_ENCAP, encapPayload))
			}

			c.egress <- encapResponse
			// c.codecComponent.Enqueue(encapResponse)
		}
	}
}

func (c TunnelComponent) ProcessIngressMessages() {
	for {
		// msg := c.codecComponent.Dequeue()
		var msg messages.MessageWrapper
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

func (c TunnelComponent) Enqueue(msg messages.MessageWrapper) {
	c.egress <- msg
}

func (c TunnelComponent) Dequeue() messages.MessageWrapper {
	msg := <-c.ingress
	return msg
}
