package gat

type Request struct {
	session *Session
	route   string
	payload []byte
	reqId   uint64
}

func (r *Request) Session() *Session {
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
