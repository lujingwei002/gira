package sproto

import (
	"context"
	"log"
	"reflect"

	"github.com/lujingwei/gira"
	gosproto "github.com/xjdrew/gosproto"
)

type Sproto struct {
	rpc *gosproto.Rpc
}

var (
	typeOfError          = reflect.TypeOf((*error)(nil)).Elem()
	typeOfContext        = reflect.TypeOf((*context.Context)(nil)).Elem()
	typeOfSprotoResponse = reflect.TypeOf((*SprotoResponse)(nil)).Elem()
)

type SprotoRequest interface {
}

type SprotoPush interface {
	GetProtoName() string
}

type SprotoResponse interface {
	SetErrorCode(v int32)
	SetErrorMsg(v string)
}

type sprotoHandlerMethod struct {
	method   reflect.Method
	typ      reflect.Type
	isRawArg bool
}

type SprotoHandler struct {
	typ     reflect.Type
	methods map[string]*sprotoHandlerMethod
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
func (self *Sproto) RequestDecode(packed []byte) (gosproto.RpcMode, string, int32, interface{}, error) {
	return self.rpc.Dispatch(packed)
}

// 编码
func (self *Sproto) RequestEncode(name string, session int32, req interface{}) (data []byte, err error) {
	return self.rpc.RequestEncode(name, session, req)
}

// 编码
func (self *Sproto) PushEncode(req SprotoPush) (data []byte, err error) {
	return self.rpc.RequestEncode(req.GetProtoName(), 0, req)
}

func (self *Sproto) ResponseEncode(name string, session int32, response interface{}) (data []byte, err error) {
	return self.rpc.ResponseEncode(name, session, response)
}

// 根据协议，调用handler的相应方法
func (self *Sproto) RpcDispatch(ctx context.Context, handler *SprotoHandler, receiver interface{}, route string, session int32, req interface{}) ([]byte, error) {
	response, _ := handler.dispatch(ctx, receiver, route, req)
	if response == nil {
		protocol := self.rpc.GetProtocolByName(route)
		response = reflect.New(protocol.Response.Elem()).Interface()
	}
	return self.rpc.ResponseEncode(route, session, response)
}

func isSprotoHandlerMethod(method reflect.Method) bool {
	mt := method.Type
	if method.PkgPath != "" {
		return false
	}
	if mt.NumIn() != 3 {
		return false
	}
	if mt.NumOut() != 2 {
		return false
	}
	if t1 := mt.In(1); t1 != typeOfContext {
		return false
	}
	if t2 := mt.In(2); t2.Kind() != reflect.Ptr {
		return false
	}
	if t0 := mt.Out(0); t0.Kind() != reflect.Ptr || !t0.Implements(typeOfSprotoResponse) {
		return false
	}
	if mt.Out(1) != typeOfError {
		return false
	}
	return true
}

func (self *SprotoHandler) suitableHandlerMethods(typ reflect.Type) map[string]*sprotoHandlerMethod {
	methods := make(map[string]*sprotoHandlerMethod)
	log.Println(typ)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mt := method.Type
		mn := method.Name
		if isSprotoHandlerMethod(method) {
			methods[mn] = &sprotoHandlerMethod{method: method, typ: mt.In(2)}
		}
	}
	return methods
}

func (self *SprotoHandler) register(handler interface{}) error {
	log.Println(handler)

	self.methods = self.suitableHandlerMethods(reflect.TypeOf(handler))
	return nil
}

func (self *SprotoHandler) dispatch(ctx context.Context, receiver interface{}, route string, r interface{}) (interface{}, error) {
	handler, found := self.methods[route]
	if !found {
		return nil, gira.ErrSprotoHandlerNotImplement
	}
	args := []reflect.Value{reflect.ValueOf(receiver), reflect.ValueOf(ctx), reflect.ValueOf(r)}
	result := handler.method.Func.Call(args)
	if err := result[1].Interface(); err == nil {
		return result[0].Interface(), nil
	} else {
		return result[0].Interface(), result[1].Interface().(error)
	}
}
