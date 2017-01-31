package codec

import "log"

import "encoding/binary"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/stack/cache"
import "github.com/chris-wood/spud/stack/pit"
import "github.com/chris-wood/spud/stack/component/connector"

const PKT_INTEREST uint8 = 0x00
const PKT_CONTENT uint8 = 0x01
const DEFAULT_HOP_LIMIT uint8 = 0xFF
const CODEC_SCHEMA_VERSION uint8 = 0x01

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
	return (uint16(bytes[0])<<8) | uint16(bytes[1])
}

func buildPacket(messageType uint16, optionalHeaderBytes, packetBytes []byte) []byte {
	header := make([]byte, 8)

	header[0] = CODEC_SCHEMA_VERSION
	switch messageType {
	case codec.T_INTEREST:
		header[1] = PKT_INTEREST
		header[4] = DEFAULT_HOP_LIMIT
	case codec.T_OBJECT:
		header[1] = PKT_CONTENT
	}

	packetLength := uint16(len(packetBytes) + len(optionalHeaderBytes) + 8)
	// TODO: check for overflow
	binary.BigEndian.PutUint16(header[2:], packetLength)

	headerLength := uint8(len(optionalHeaderBytes) + 8)
	// TODO: check for overflow
	header[7] = headerLength

	wireFormat := append(header, optionalHeaderBytes...)
	wireFormat = append(wireFormat, packetBytes...)
	return wireFormat
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
		wireFormat := buildPacket(messageType, optionalHeader, messageBytes)

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

func (c CodecComponent) Enqueue(msg *messages.MessageWrapper) {
	c.egress <- msg
}

func (c CodecComponent) Dequeue() *messages.MessageWrapper {
	msg := <-c.ingress
	return msg
}
