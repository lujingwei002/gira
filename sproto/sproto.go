package sproto

import (
	"context"
	"reflect"

	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
	gosproto "github.com/xjdrew/gosproto"
)

type Sproto struct {
	rpc *gosproto.Rpc
}

var (
	typeOfError          = reflect.TypeOf((*error)(nil)).Elem()
	typeOfContext        = reflect.TypeOf((*context.Context)(nil)).Elem()
	typeOfSprotoRequest  = reflect.TypeOf((*SprotoRequest)(nil)).Elem()
	typeOfSprotoResponse = reflect.TypeOf((*SprotoResponse)(nil)).Elem()
	typeOfSprotoPush     = reflect.TypeOf((*SprotoPush)(nil)).Elem()
	typeOfSprotoPushArr  = reflect.TypeOf(([]SprotoPush)(nil))
)

type SprotoRequest interface {
	GetRequestName() string
}

type SprotoPush interface {
	GetPushName() string
}

type SprotoResponse interface {
	SetErrorCode(v int32)
	SetErrorMsg(v string)
}

type proto_type int

const (
	proto_handler_type_request = 1
	proto_handler_type_push    = 2
)

func (self proto_type) String() string {
	switch self {
	case proto_handler_type_push:
		return "Push"
	case proto_handler_type_request:
		return "Request"
	default:
		return "unsupport"
	}
}

type sproto_handler_method struct {
	method    reflect.Method
	protoType proto_type
}

type SprotoHandler struct {
	typ     reflect.Type
	methods map[string]*sproto_handler_method
}

func RegisterRpc(protocols []*gosproto.Protocol) (*Sproto, error) {
	if rpc, err := gosproto.NewRpc(protocols); err != nil {
		return nil, err
	} else {
		self := &Sproto{
			rpc: rpc,
		}
		return self, nil
	}
}

// 生成handler
func (self *Sproto) RegisterHandler(handler interface{}) *SprotoHandler {
	handlers := &SprotoHandler{}
	if err := handlers.register(handler); err != nil {
		return nil
	}
	return handlers
}

// 解码
func (self *Sproto) RequestDecode(packed []byte) (mode gosproto.RpcMode, route string, session int32, resp SprotoRequest, err error) {
	var sp interface{}
	mode, route, session, sp, err = self.rpc.Dispatch(packed)
	if err != nil {
		return
	}
	resp = sp.(SprotoRequest)
	return
}

func (self *Sproto) ResponseDecode(packed []byte) (mode gosproto.RpcMode, route string, session int32, resp SprotoResponse, err error) {
	var sp interface{}
	mode, route, session, sp, err = self.rpc.Dispatch(packed)
	if err != nil {
		return
	}
	resp = sp.(SprotoResponse)
	return
}

func (self *Sproto) PushDecode(packed []byte) (mode gosproto.RpcMode, route string, session int32, resp SprotoPush, err error) {
	var sp interface{}
	mode, route, session, sp, err = self.rpc.Dispatch(packed)
	if err != nil {
		return
	}
	resp = sp.(SprotoPush)
	return
}

func (self *Sproto) NewResponse(req SprotoRequest) (resp SprotoResponse, err error) {
	proto := self.rpc.GetProtocolByName(req.GetRequestName())
	if proto == nil {
		err = gira.ErrTodo
		return
	}
	sp := reflect.New(proto.Response.Elem()).Interface()
	if sp == nil {
		err = gira.ErrTodo
		return
	}
	resp = sp.(SprotoResponse)
	return
}

// 编码
func (self *Sproto) RequestEncode(name string, session int32, req interface{}) (data []byte, err error) {
	return self.rpc.RequestEncode(name, session, req)
}

// 编码
func (self *Sproto) PushEncode(req SprotoPush) (data []byte, err error) {
	return self.rpc.RequestEncode(req.GetPushName(), 0, req)
}

func (self *Sproto) ResponseEncode(name string, session int32, response interface{}) (data []byte, err error) {
	return self.rpc.ResponseEncode(name, session, response)
}

func (self *Sproto) StructEncode(req interface{}) (data []byte, err error) {
	return gosproto.Encode(req)
}

func (self *Sproto) StructDecode(data []byte, req interface{}) error {
	if _, err := gosproto.Decode(data, req); err != nil {
		return err
	} else {
		return nil
	}
}

// 根据协议，调用handler的相应方法
func (self *Sproto) RequestDispatch(ctx context.Context, handler *SprotoHandler, receiver interface{}, route string, session int32, req interface{}) (dataResp []byte, dataPushArr [][]byte, err error) {
	resp, pushArr, err := handler.RequestDispatch(ctx, receiver, route, req)
	// response
	if resp == nil {
		protocol := self.rpc.GetProtocolByName(route)
		elem := reflect.New(protocol.Response.Elem())
		if !elem.IsNil() {
			resp = elem.Interface()
		}
		if resp == nil {
			err = gira.ErrSprotoResponseConversion
			return
		}
	}
	proto, ok := resp.(SprotoResponse)
	if !ok {
		err = gira.ErrSprotoResponseConversion
		return
	}
	if proto != nil && err != nil {
		proto.SetErrorCode(gira.ErrCode(err))
		proto.SetErrorMsg(gira.ErrMsg(err))
	}
	dataResp, err = self.rpc.ResponseEncode(route, session, proto)
	if err != nil {
		return
	}
	// push
	if pushArr != nil {
		var dataPush []byte
		for _, v := range pushArr {
			proto, ok := v.(SprotoPush)
			if !ok {
				err = gira.ErrSprotoPushConversion
				return
			}
			dataPush, err = self.rpc.RequestEncode(proto.GetPushName(), 0, proto)
			if err != nil {
				return
			}
			dataPushArr = append(dataPushArr, dataPush)
		}
	}
	return
}

// 根据协议，调用handler的相应方法
func (self *Sproto) PushDispatch(ctx context.Context, handler *SprotoHandler, receiver interface{}, route string, req SprotoPush) error {
	return handler.PushDispatch(ctx, receiver, route, req)
}

func isSprotoHandlerMethod(method reflect.Method) *sproto_handler_method {
	mt := method.Type
	if method.PkgPath != "" {
		return nil
	}
	if mt.NumIn() != 3 {
		return nil
	}
	if arg1 := mt.In(1); arg1 != typeOfContext {
		return nil
	}
	if arg2 := mt.In(2); arg2.Kind() != reflect.Ptr {
		return nil
	}
	if arg2 := mt.In(2); arg2.Implements(typeOfSprotoRequest) {
		if mt.NumOut() != 2 && mt.NumOut() != 3 {
			return nil
		}
		if r0 := mt.Out(0); r0.Kind() != reflect.Ptr || !r0.Implements(typeOfSprotoResponse) {
			return nil
		}
		if mt.NumOut() == 2 {
			if mt.Out(1) != typeOfError {
				return nil
			}
		} else if mt.NumOut() == 3 {
			if mt.Out(1) != typeOfSprotoPushArr {
				return nil
			}
			if mt.Out(2) != typeOfError {
				return nil
			}
		}
		return &sproto_handler_method{method: method, protoType: proto_handler_type_request}
	} else if arg2 := mt.In(2); arg2.Implements(typeOfSprotoPush) {
		if mt.NumOut() != 1 {
			return nil
		}
		if mt.Out(0) != typeOfError {
			return nil
		}
		return &sproto_handler_method{method: method, protoType: proto_handler_type_push}
	}

	return nil
}

func (self *SprotoHandler) suitableHandlerMethods(typ reflect.Type) map[string]*sproto_handler_method {
	methods := make(map[string]*sproto_handler_method)
	log.Info(typ)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mn := method.Name
		if handler := isSprotoHandlerMethod(method); handler != nil {
			log.Infow("registr handler", "type", handler.protoType, "name", mn)
			methods[mn] = handler
		}
	}
	return methods
}

func (self *SprotoHandler) register(handler interface{}) error {
	self.methods = self.suitableHandlerMethods(reflect.TypeOf(handler))
	return nil
}

func (self *SprotoHandler) HasRoute(route string) bool {
	_, found := self.methods[route]
	if !found {
		return false
	}
	return true
}

func (self *SprotoHandler) RequestDispatch(ctx context.Context, receiver interface{}, route string, r interface{}) (resp interface{}, push []SprotoPush, err error) {
	handler, found := self.methods[route]
	if !found {
		log.Warnw("sproto handler not found", "name", route)
		err = gira.ErrSprotoHandlerNotImplement
		return
	}
	args := []reflect.Value{reflect.ValueOf(receiver), reflect.ValueOf(ctx), reflect.ValueOf(r)}
	result := handler.method.Func.Call(args)
	if handler.method.Type.NumOut() == 3 {
		result0 := result[0]
		result1 := result[1]
		result2 := result[2]
		if !result0.IsNil() {
			resp = result0.Interface()
		}
		if !result2.IsNil() {
			err = result2.Interface().(error)
		}
		if !result1.IsNil() {
			push = (result1.Interface()).([]SprotoPush)
		}
	} else if handler.method.Type.NumOut() == 2 {
		result0 := result[0]
		result1 := result[1]
		if !result0.IsNil() {
			resp = result0.Interface()
		}
		if !result1.IsNil() {
			err = result1.Interface().(error)
		}
	}
	return
}

func (self *SprotoHandler) PushDispatch(ctx context.Context, receiver interface{}, route string, r interface{}) (err error) {
	handler, found := self.methods[route]
	if !found {
		log.Warnw("sproto handler not found", "name", route)
		err = gira.ErrSprotoHandlerNotImplement
		return
	}
	args := []reflect.Value{reflect.ValueOf(receiver), reflect.ValueOf(ctx), reflect.ValueOf(r)}
	result := handler.method.Func.Call(args)
	if handler.method.Type.NumOut() == 1 {
		result0 := result[0]
		if !result0.IsNil() {
			err = result0.Interface().(error)
		}
	}
	return
}
