package codec

import "encoding/binary"
import "github.com/chris-wood/spud/messages"
import "github.com/chris-wood/spud/codec"
import "github.com/chris-wood/spud/stack/component/connector"

const PKT_INTEREST uint8 = 0x00
const PKT_CONTENT uint8 = 0x01
const DEFAULT_HOP_LIMIT uint8 = 0xFF
const CODEC_SCHEMA_VERSION uint8 = 0x01

type CodecComponent struct {
    ingress chan messages.Message
    egress chan messages.Message

    connector connector.ForwarderConnector
}

func NewCodecComponent(conn connector.ForwarderConnector) CodecComponent {
    egress := make(chan messages.Message)
    ingress := make(chan messages.Message)

    return CodecComponent{ingress: ingress, egress: egress, connector: conn}
}

func readWord(bytes []byte) uint16 {
    return uint16(bytes[0] << 8 | bytes[1])
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
    for ;; {
        msg := <- c.egress

        // 1. Encode the message
        messageBytes := msg.Encode()

        // 2. Encode optional headers, if present
        // XX: currently not implemented
        optionalHeader := make([]byte, 0)

        // 3. Prepend the fixed header and make the final packet
        messageType := readWord(messageBytes)
        wireFormat := buildPacket(messageType, optionalHeader, messageBytes)

        // 4. Send the wiret format packet to the forwarder connector
        c.connector.Write(wireFormat)
    }
}

func (c CodecComponent) ProcessIngressMessages() {
    decoder := codec.Decoder{}
    for ;; {
        msgBytes := c.connector.Read()

        // 1. Extract the message bytes (strip headers)
        // packetLength := readWord(msgBytes[2:4])
        headerLength := msgBytes[7]

        // 2. Decode the message
        decodedTlV := decoder.Decode(msgBytes[headerLength:])
        message, err := messages.CreateFromTLV(decodedTlV)

        // 3. Enqueue in the upstream (ingress) queue
        if err == nil {
            c.ingress <- message
        }
    }
}

func (c CodecComponent) Enqueue(msg messages.Message) {
    c.egress <- msg
}

func (c CodecComponent) Dequeue() messages.Message {
    msg := <-c.ingress
    return msg
}
