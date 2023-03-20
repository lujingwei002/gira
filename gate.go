package gira

import "context"

type GateConn interface {
	ID() uint64
	Close() error
	Kick(reason string)
	Recv(ctx context.Context) (GateRequest, error)
	Push(route string, data []byte) error
	Response(mid uint64, data []byte) error
}

type GateClient interface {
	Recv(ctx context.Context) (typ int, route string, reqId uint64, data []byte, err error)
	Notify(route string, data []byte) error
	Request(route string, reqId uint64, data []byte) error
	Close() error
}

type GateHandler interface {
	OnGateStream(conn GateConn)
}

type GateRequest interface {
	Response(data []byte) error
	Payload() []byte
	ReqId() uint64
}

const (
	GateMessageType_RESPONSE = 2
	GateMessageType_PUSH     = 3
)
