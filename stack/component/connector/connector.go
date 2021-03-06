package connector

import "net"

// The generic forwarder connector interface
// A forwarder connector is used to read raw packets from and write raw packets
//   to the forwarder.
type ForwarderConnector interface {
	Read() []byte
	Write([]byte)
}

type connectorError struct {
	problem string
}

func (ce connectorError) Error() string {
	return ce.problem
}

// Loopback forward connector
type LoopbackForwarderConnector struct {
	buffer chan []byte
}

func NewLoopbackForwarderConnector() (*LoopbackForwarderConnector, error) {
	return &LoopbackForwarderConnector{buffer: make(chan []byte)}, nil
}

func (fc LoopbackForwarderConnector) Read() []byte {
	slice := <-fc.buffer
	return slice
}

func (fc LoopbackForwarderConnector) Write(bytes []byte) {
	fc.buffer <- bytes
}

// The TCP forwarder connector
type TCPForwarderConnector struct {
	conn   net.Conn
	buffer []byte
}

func NewAthenaTCPForwarderConnector(location string) (*TCPForwarderConnector, error) {
	connection, err := net.Dial("tcp", location)
	if err != nil {
		return nil, connectorError{"Unable to connect to the forwarder at " + location}
	}

	return &TCPForwarderConnector{conn: connection, buffer: make([]byte, 64000)}, nil
}

func (fc TCPForwarderConnector) Read() []byte {
	num, _ := fc.conn.Read(fc.buffer)
	return fc.buffer[:num]
}

func (fc TCPForwarderConnector) Write(bytes []byte) {
	fc.conn.Write(bytes)
	// XX: this ignores the return result. what would we do if it failed?
}
