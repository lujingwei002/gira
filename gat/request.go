package gat

import "github.com/lujingwei/gira"

type Request struct {
	session *Session
	route   string
	payload []byte
	reqId   uint64
}

func (r *Request) Session() gira.GateConn {
	return r.session
}
func (r *Request) Payload() []byte {
	return r.payload
}

func (r *Request) ReqId() uint64 {
	return r.reqId
}

func (r *Request) Response(data []byte) error {
	return r.session.Response(r.reqId, data)
}

func (r *Request) Push(route string, data []byte) error {
	return r.session.Push(route, data)
}
