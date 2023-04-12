package packet

import (
	"errors"
	"fmt"
)

// Type represents the network packet's type such as: handshake and so on.
type Type byte

const (
	_       Type = iota
	TypeMin      = Handshake
	// Handshake represents a handshake: request(client) <====> handshake response(server)
	Handshake = 0x01
	// HandshakeAck represents a handshake ack from client to server
	HandshakeAck = 0x02
	// Data represents a common data packet
	Data = 0x03
	// Heartbeat represents a heartbeat
	Heartbeat = 0x04
	// Kick represents a kick off packet
	Kick           = 0x05 // disconnect message from server
	ServerDown     = 0x06
	ServerMaintain = 0x07
	UserInstead    = 0x08
	ServerSuspend  = 0x09
	ServerResume   = 0x0a
	TypeMax        = ServerResume
)

// ErrWrongPacketType represents a wrong packet type.
var ErrWrongPacketType = errors.New("wrong packet type")

// Packet represents a network packet.
type Packet struct {
	Type   Type
	Length int
	Data   []byte
}

// New create a Packet instance.
func NewPacket() *Packet {
	return &Packet{}
}

// String represents the Packet's in text mode.
func (p *Packet) String() string {
	return fmt.Sprintf("Type: %d, Length: %d, Data: %s", p.Type, p.Length, string(p.Data))
}
