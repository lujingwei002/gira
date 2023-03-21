package sproto

import (
	"context"
	"reflect"

	"github.com/lujingwei/gira/log"

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
	typeOfSprotoPushArr  = reflect.TypeOf(([]SprotoPush)(nil))
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
	method       reflect.Method
	isReturnPush bool
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
func (self *Sproto) RpcDispatch(ctx context.Context, handler *SprotoHandler, receiver interface{}, route string, session int32, req interface{}) (dataResp []byte, dataPushArr [][]byte, err error) {
	resp, pushArr, err := handler.dispatch(ctx, receiver, route, req)
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
			dataPush, err = self.rpc.RequestEncode(proto.GetProtoName(), 0, proto)
			if err != nil {
				return
			}
			dataPushArr = append(dataPushArr, dataPush)
		}
	}
	return
}

func isSprotoHandlerMethod(method reflect.Method) bool {
	mt := method.Type
	if method.PkgPath != "" {
		return false
	}
	if mt.NumIn() != 3 {
		return false
	}
	if mt.NumOut() != 2 && mt.NumOut() != 3 {
		return false
	}
	if arg1 := mt.In(1); arg1 != typeOfContext {
		return false
	}
	if arg2 := mt.In(2); arg2.Kind() != reflect.Ptr {
		return false
	}
	if r0 := mt.Out(0); r0.Kind() != reflect.Ptr || !r0.Implements(typeOfSprotoResponse) {
		return false
	}
	if mt.NumOut() == 2 {
		if mt.Out(1) != typeOfError {
			return false
		}
	} else if mt.NumOut() == 3 {
		if mt.Out(1) != typeOfSprotoPushArr {
			return false
		}
		if mt.Out(2) != typeOfError {
			return false
		}
	}
	return true
}

func (self *SprotoHandler) suitableHandlerMethods(typ reflect.Type) map[string]*sprotoHandlerMethod {
	methods := make(map[string]*sprotoHandlerMethod)
	log.Info(typ)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mn := method.Name
		if isSprotoHandlerMethod(method) {
			gira.Infow("registr handler", "name", mn)
			handler := &sprotoHandlerMethod{method: method}
			if method.Type.NumOut() == 3 {
				handler.isReturnPush = true
			}
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

func (self *SprotoHandler) dispatch(ctx context.Context, receiver interface{}, route string, r interface{}) (resp interface{}, push []SprotoPush, err error) {
	handler, found := self.methods[route]
	if !found {
		gira.Infow("sproto handler not found", "name", route)
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
	} else {
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
