package transport

import "log"
import "time"

import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/stack/component"

type TimeoutEvent struct {
	timer     chan bool
	cancelled bool
}

func Start(identity string, duration time.Duration, callback func(identity string)) *TimeoutEvent {
	timer := make(chan bool, 1)
	event := &TimeoutEvent{timer, false}

	go func(timeoutEvent *TimeoutEvent) {
		select {
		case <-timeoutEvent.timer:
			timeoutEvent.cancelled = true
		case <-time.After(duration):
			callback(identity)
		}
	}(event)

	return event
}

func (t *TimeoutEvent) Cancel() bool {
	if t.cancelled {
		return false
	} else {
		t.timer <- true
		return true
	}
}

type TransportComponent struct {
	ingress    chan *messages.MessageWrapper
	egress     chan *messages.MessageWrapper
	downstream component.Component

	pendingTable map[string]*TimeoutEvent
}

func NewTransportComponent(downstream component.Component) *TransportComponent {
	egress := make(chan *messages.MessageWrapper)
	ingress := make(chan *messages.MessageWrapper)

	log.Println("Created Transport component")

	return &TransportComponent{
		ingress:      ingress,
		egress:       egress,
		downstream:   downstream,
		pendingTable: make(map[string]*TimeoutEvent),
	}
}

func (c *TransportComponent) HandleTimeout(identity string) {
	// XXX
	log.Println("XXX: timeout handler not implemented")
}

func (c *TransportComponent) ProcessEgressMessages() {
	for {
		msg := <-c.egress

		// XXX: Do transport stuff

		// Start a timeout if it's a request
		if msg.GetPacketType() == codec.T_INTEREST {
			identifier := msg.Identifier()
			c.pendingTable[identifier] = Start(identifier, time.Second, c.HandleTimeout)
		}

		// Send the message downstream
		c.downstream.Push(msg)
	}
}

func (c *TransportComponent) ProcessIngressMessages() {
	for {
		msg := c.downstream.Pop()

		// Cancel timeout
		if msg.GetPacketType() != codec.T_INTEREST {
			identifier := msg.Identifier()
			if match, ok := c.pendingTable[identifier]; ok {
				if cancelled := match.Cancel(); cancelled {

					// XXX: Do transport stuff

					delete(c.pendingTable, identifier)

					// Pass upstream
					c.ingress <- msg
				}
			} else {
				log.Println("Unknown message received. Dropping.")
			}
		} else {
			// All interests automatically go up the pipe
			c.ingress <- msg
		}
	}
}

func (c TransportComponent) Push(msg *messages.MessageWrapper) {
	c.egress <- msg
}

func (c TransportComponent) Pop() *messages.MessageWrapper {
	msg := <-c.ingress
	return msg
}
