package sproto

//
// proto接口的sproto实现
// 包括消息解码，编码和消息路由功能
//

import (
	"context"
	"reflect"

	"github.com/lujingwei002/gira/log"

	"github.com/lujingwei002/gira"
	gosproto "github.com/xjdrew/gosproto"
)

type sproto struct {
	rpc *gosproto.Rpc
}

var (
	typeOfError          = reflect.TypeOf((*error)(nil)).Elem()
	typeOfContext        = reflect.TypeOf((*context.Context)(nil)).Elem()
	typeOfSprotoRequest  = reflect.TypeOf((*gira.ProtoRequest)(nil)).Elem()
	typeOfSprotoResponse = reflect.TypeOf((*gira.ProtoResponse)(nil)).Elem()
	typeOfSprotoPush     = reflect.TypeOf((*gira.ProtoPush)(nil)).Elem()
	typeOfSprotoPushArr  = reflect.TypeOf(([]gira.ProtoPush)(nil))
)

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

type sproto_handler struct {
	typ     reflect.Type
	methods map[string]*sproto_handler_method
}

func RegisterRpc(protocols []*gosproto.Protocol) (*sproto, error) {
	if rpc, err := gosproto.NewRpc(protocols); err != nil {
		return nil, err
	} else {
		self := &sproto{
			rpc: rpc,
		}
		return self, nil
	}
}

// 生成handler
func (self *sproto) RegisterHandler(handler interface{}) gira.ProtoHandler {
	handlers := &sproto_handler{}
	if err := handlers.register(handler); err != nil {
		return nil
	}
	return handlers
}

// request解码
func (self *sproto) RequestDecode(packed []byte) (route string, session int32, resp gira.ProtoRequest, err error) {
	var sp interface{}
	_, route, session, sp, err = self.rpc.Dispatch(packed)
	if err != nil {
		return
	}
	resp = sp.(gira.ProtoRequest)
	return
}

// response解码
func (self *sproto) ResponseDecode(packed []byte) (route string, session int32, resp gira.ProtoResponse, err error) {
	var sp interface{}
	_, route, session, sp, err = self.rpc.Dispatch(packed)
	if err != nil {
		return
	}
	resp = sp.(gira.ProtoResponse)
	return
}

// push解码
func (self *sproto) PushDecode(packed []byte) (route string, session int32, resp gira.ProtoPush, err error) {
	var sp interface{}
	_, route, session, sp, err = self.rpc.Dispatch(packed)
	if err != nil {
		return
	}
	resp = sp.(gira.ProtoPush)
	return
}

// 从request创建response
func (self *sproto) NewResponse(req gira.ProtoRequest) (resp gira.ProtoResponse, err error) {
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
	resp = sp.(gira.ProtoResponse)
	return
}

// request编码
func (self *sproto) RequestEncode(name string, session int32, req interface{}) (data []byte, err error) {
	return self.rpc.RequestEncode(name, session, req)
}

// push编码
func (self *sproto) PushEncode(req gira.ProtoPush) (data []byte, err error) {
	return self.rpc.RequestEncode(req.GetPushName(), 0, req)
}

// response编码
func (self *sproto) ResponseEncode(name string, session int32, response interface{}) (data []byte, err error) {
	return self.rpc.ResponseEncode(name, session, response)
}

// struct编码
func (self *sproto) StructEncode(req interface{}) (data []byte, err error) {
	return gosproto.Encode(req)
}

// struct解码
func (self *sproto) StructDecode(data []byte, req interface{}) error {
	if _, err := gosproto.Decode(data, req); err != nil {
		return err
	} else {
		return nil
	}
}

// 根据协议，调用handler的相应方法
func (self *sproto) RequestDispatch(ctx context.Context, handler gira.ProtoHandler, receiver interface{}, route string, session int32, req interface{}) (dataResp []byte, pushArr []gira.ProtoPush, err error) {
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
	proto, ok := resp.(gira.ProtoResponse)
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
	// if pushArr != nil {
	// 	dataPushArr = make([][]byte, 0)
	// 	for _, v := range pushArr {
	// 		proto, ok := v.(gira.ProtoPush)
	// 		if !ok {
	// 			err = gira.ErrSprotoPushConversion
	// 			return
	// 		}
	// 		var dataPush []byte
	// 		dataPush, err = self.rpc.RequestEncode(proto.GetPushName(), 0, proto)
	// 		if err != nil {
	// 			return
	// 		}
	// 		dataPushArr = append(dataPushArr, dataPush)
	// 	}
	// }
	return
}

// 根据协议，调用handler的相应方法
func (self *sproto) PushDispatch(ctx context.Context, handler gira.ProtoHandler, receiver interface{}, route string, req gira.ProtoPush) error {
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

func (self *sproto_handler) suitableHandlerMethods(typ reflect.Type) map[string]*sproto_handler_method {
	methods := make(map[string]*sproto_handler_method)
	log.Info(typ)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mn := method.Name
		if handler := isSprotoHandlerMethod(method); handler != nil {
			log.Infow("registr sproto handler", "type", handler.protoType, "name", mn)
			methods[mn] = handler
		}
	}
	return methods
}

func (self *sproto_handler) register(handler interface{}) error {
	self.methods = self.suitableHandlerMethods(reflect.TypeOf(handler))
	return nil
}

// 是否存在路由
func (self *sproto_handler) HasRoute(route string) bool {
	_, found := self.methods[route]
	if !found {
		return false
	}
	return true
}

// 处理request
func (self *sproto_handler) RequestDispatch(ctx context.Context, receiver interface{}, route string, r interface{}) (resp interface{}, push []gira.ProtoPush, err error) {
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
			push = (result1.Interface()).([]gira.ProtoPush)
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

// 处理push
func (self *sproto_handler) PushDispatch(ctx context.Context, receiver interface{}, route string, r interface{}) (err error) {
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
