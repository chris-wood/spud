package codec

import (
	"github.com/chris-wood/spud/stack/component/connector"
	"github.com/chris-wood/spud/messages"
	"github.com/chris-wood/spud/stack/cache"
	"github.com/chris-wood/spud/stack/pit"
	"github.com/chris-wood/spud/codec"
	"encoding/binary"
	"fmt"
	"math"
	"log"
)

const PKT_INTEREST uint8 = 0x00
const PKT_CONTENT uint8 = 0x01
const DEFAULT_HOP_LIMIT uint8 = 0xFF
const CODEC_SCHEMA_VERSION uint8 = 0x01

type codecError struct {
	prob string
}

func (e codecError) Error() string {
	return fmt.Sprintf("%s", e.prob)
}

type CodecComponent struct {
	ingress chan *messages.MessageWrapper
	egress  chan *messages.MessageWrapper

	stackCache *cache.Cache
	stackPit   *pit.PIT
	connector  connector.ForwarderConnector
}

func NewCodecComponent(conn connector.ForwarderConnector, stackCache *cache.Cache, stackPit *pit.PIT) CodecComponent {
	egress := make(chan *messages.MessageWrapper)
	ingress := make(chan *messages.MessageWrapper)

	return CodecComponent{ingress: ingress, egress: egress, connector: conn, stackCache: stackCache, stackPit: stackPit}
}

func readWord(bytes []byte) uint16 {
	return (uint16(bytes[0]) << 8) | uint16(bytes[1])
}

func buildPacket(messageType uint16, optionalHeaderBytes, packetBytes []byte) ([]byte, error) {
	header := make([]byte, 8)

	header[0] = CODEC_SCHEMA_VERSION
	switch messageType {
	case codec.T_INTEREST:
		header[1] = PKT_INTEREST
		header[4] = DEFAULT_HOP_LIMIT
	case codec.T_OBJECT:
		header[1] = PKT_CONTENT
	}

	// Check for overflow
	if (len(packetBytes) + len(optionalHeaderBytes) + 8 > math.MaxUint16) {
		return nil, codecError{"Packet length exceeded the wire format capacity"}
	}

	packetLength := uint16(len(packetBytes) + len(optionalHeaderBytes) + 8)
	binary.BigEndian.PutUint16(header[2:], packetLength)

	// Check for header length overflow
	if (len(optionalHeaderBytes) + 8 > math.MaxUint8) {
		return nil, codecError{"Header length exceeded a single octet"}
	}
	headerLength := uint8(len(optionalHeaderBytes) + 8)
	header[7] = headerLength

	wireFormat := append(header, optionalHeaderBytes...)
	wireFormat = append(wireFormat, packetBytes...)

	return wireFormat, nil
}

func (c CodecComponent) ProcessEgressMessages() {
	for {
		msg := <-c.egress

		// Encode the message
		messageBytes := msg.Encode()

		// Encode optional headers, if present
		// XXX: currently not implemented
		optionalHeader := make([]byte, 0)

		// Prepend the fixed header and make the final packet
		messageType := readWord(messageBytes)
		wireFormat, err := buildPacket(messageType, optionalHeader, messageBytes)
		if err != nil {
			log.Println(err)
			continue
		}

		// If we have a content object, insert it into the cache and forward if a PIT entry awaits
		// Otherwise, if it's an interest without a PIT entry, forward it
		// Otherwise, it's an interest with a PIT entry, so aggregate
		if messageType != codec.T_INTEREST {
			c.stackCache.Insert(msg.Identifier(), wireFormat)
			if _, found := c.stackPit.Lookup(msg.Identifier()); found {
				c.connector.Write(wireFormat)
			}
		} else if _, found := c.stackPit.Lookup(msg.Identifier()); !found {
			c.connector.Write(wireFormat)
		} else {
			c.stackPit.Insert(msg.Identifier(), msg) // this should aggregate
		}
	}
}

func (c CodecComponent) ProcessIngressMessages() {
	decoder := codec.Decoder{}
	for {
		msgBytes := c.connector.Read()
		if len(msgBytes) > 8 {
			// Extract the message bytes (strip headers)
			packetLength := readWord(msgBytes[2:4])
			headerLength := msgBytes[7]

			// Ensure the packet length is correct
			if len(msgBytes) != int(packetLength) {
				log.Printf("Packet length mismatch. Expected %d, got %d. Dropping. %s", int(packetLength), len(msgBytes))
				continue
			}

			// Decode the message (skipping past the packet header)
			decodedTlV := decoder.Decode(msgBytes[headerLength:])
			message, err := messages.CreateFromTLV(decodedTlV)
			if err == nil {
				identity := message.Identifier()

				// If the response is cached, just serve it
				match, isPresent := c.stackCache.Lookup(identity) // XXX: not necessary for responses
				if isPresent && message.GetPacketType() == codec.T_INTEREST {
					c.connector.Write(match)
					continue
				}

				// If the response is not cached, but it's in the PIT, drop it
				_, inPit := c.stackPit.Lookup(identity)
				if inPit && message.GetPacketType() == codec.T_INTEREST {
					log.Println("Supressing duplicate PIT entry", msgBytes)
					continue
				}

				// Otherwise, forward it up and insert it into the PIT
				c.stackPit.Insert(message.Identifier(), message)
				c.ingress <- message
			} else {
				log.Println("Failed to decode the message", err)
			}
		}
	}
}

func (c CodecComponent) Push(msg *messages.MessageWrapper) {
	c.egress <- msg
}

func (c CodecComponent) Pop() *messages.MessageWrapper {
	msg := <-c.ingress
	return msg
}
