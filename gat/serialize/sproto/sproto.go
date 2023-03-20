package sproto

import (
	"errors"

	sproto "github.com/xjdrew/gosproto"
)

type Serializer struct {
	rpc *sproto.Rpc
}

type Request struct {
	Name      string
	SessionId int32
	Data      interface{}
}

type Response struct {
	Name      string
	SessionId int32
	Mode      sproto.RpcMode
	Data      interface{}
}

func NewSerializer(rpc *sproto.Rpc) *Serializer {
	return &Serializer{
		rpc: rpc,
	}
}

func (s *Serializer) Marshal(v interface{}) ([]byte, error) {
	req, ok := v.(*Request)
	if !ok {
		return nil, errors.New("sproto request type error")
	}
	return s.rpc.RequestEncode(req.Name, req.SessionId, req.Data)
}

func (s *Serializer) Unmarshal(data []byte, v interface{}) error {
	resp, ok := v.(*Response)
	if !ok {
		return errors.New("sproto request type error")
	}
	mode, name, session, pack, err := s.rpc.Dispatch(data)
	if err != nil {
		return err
	}
	resp.Name = name
	resp.SessionId = session
	resp.Mode = mode
	resp.Data = pack
	return nil
}
