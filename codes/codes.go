package codes

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
)

type Logger interface {
	Errorw(msg string, keysAndValues ...interface{})
}

var defaultLogger Logger

func SetLogger(l Logger) {
	defaultLogger = l
}

// 错误码常量
const (
	CodeOk                             = 0
	CodeUnknown                        = 1
	CodeNullPointer                    = 2
	CodeNullObject                     = 3
	CodeInvalidSdkToken                = 4
	CodeInvalidArgs                    = 5
	CodeInvalidJwt                     = 6
	CodeJwtExpire                      = 7
	CodeToDo                           = 8
	CodeInvalidPassword                = 9
	CodePayOrderStatusInvalid          = 10
	CodeAccountPlatformNotSupport      = 11
	CodePayOrderAccountPlatformInvalid = 12
	CodePayOrderAmountInvalid          = 13
)
const (
	msg_OK                             = "成功"
	msg_Unknown                        = "未知错误"
	msg_NullObject                     = "null object"
	msg_InvalidSdkToken                = "无效的sdk token"
	msg_InvalidArgs                    = "无效参数"
	msg_InvalidJwt                     = "无效的token"
	msg_JwtExpire                      = "token已过期"
	msg_ToDo                           = "TODO"
	msg_InvalidPassword                = "inivalid password"
	msg_PayOrderStatusInvalid          = "pay order status invalid"
	msg_AccountPlatformNotSupport      = "account platform not support"
	msg_PayOrderAccountPlatformInvalid = "pay order account platform invalid"
	msg_PayOrderAmountInvalid          = "pay order amount invalid"
)

var (
	ErrUnknown                        = New(CodeUnknown, msg_Unknown)
	ErrNullObject                     = New(CodeNullObject, msg_NullObject)
	ErrInvalidSdkToken                = New(CodeInvalidSdkToken, msg_InvalidSdkToken)
	ErrInvalidArgs                    = New(CodeInvalidArgs, msg_InvalidArgs)
	ErrInvalidJwt                     = New(CodeInvalidJwt, msg_InvalidJwt)
	ErrJwtExpire                      = New(CodeJwtExpire, msg_JwtExpire)
	ErrToDo                           = New(CodeToDo, msg_ToDo)
	ErrInvalidPassword                = New(CodeInvalidPassword, msg_InvalidPassword)
	ErrPayOrderStatusInvalid          = New(CodePayOrderStatusInvalid, msg_PayOrderStatusInvalid)
	ErrAccountPlatformNotSupport      = New(CodeAccountPlatformNotSupport, msg_AccountPlatformNotSupport)
	ErrPayOrderAccountPlatformInvalid = New(CodePayOrderAccountPlatformInvalid, msg_PayOrderAccountPlatformInvalid)
	ErrPayOrderAmountInvalid          = New(CodePayOrderAmountInvalid, msg_PayOrderAmountInvalid)
)

func TraceErrTodo(values ...interface{}) *TraceError {
	return ErrToDo.Trace(values...)
}

// 根据code, msg, values创建error_code
func Trace(code int32, msg string, values ...interface{}) *CodeError {
	var kvs map[string]interface{}
	if len(values)%2 != 0 {
	} else if len(values) == 0 {
	} else {
		kvs = make(map[string]interface{})

		for i := 0; i < len(values); i += 2 {
			j := i + 1
			if k, ok := values[i].(string); ok {
				kvs[k] = values[j]
			}
		}
	}
	if defaultLogger != nil {
		vs := append(values, "stack", string(debug.Stack()))
		defaultLogger.Errorw(msg, vs...)
	}
	return &CodeError{
		Code:   code,
		Msg:    msg,
		Values: kvs,
	}
}

func Unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}

func Is(err error, target *CodeError) bool {
	if target == nil {
		return err == target
	}
	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable && err == target {
			return true
		}
		if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
			return true
		}
		// TODO: consider supporting target.Is(err). This would allow
		// user-definable predicates, but also may allow for coping with sloppy
		// APIs, thereby making it easier to get away with them.
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

// 根据code, msg, values创建error_code
func New(code int32, msg string, values ...interface{}) *CodeError {
	var kvs map[string]interface{}
	if len(values)%2 != 0 {
	} else if len(values) == 0 {
	} else {
		kvs = make(map[string]interface{})

		for i := 0; i < len(values); i += 2 {
			j := i + 1
			if k, ok := values[i].(string); ok {
				kvs[k] = values[j]
			}
		}
	}
	return &CodeError{
		Code:   code,
		Msg:    msg,
		Values: kvs,
	}
}

// 转换或者创建error
func NewOrCastError(err error) *CodeError {
	if err == nil {
		return nil
	}
	if e, ok := err.(*CodeError); ok {
		return e
	} else {
		msg := err.Error()
		arr := strings.Split(msg, ":")
		if len(arr) < 2 {
			return &CodeError{
				Code: CodeUnknown,
				Msg:  msg,
			}
		} else if v, err := strconv.Atoi(arr[0]); err != nil {
			return &CodeError{
				Code: CodeUnknown,
				Msg:  msg,
			}
		} else {
			return &CodeError{
				Code: int32(v),
				Msg:  arr[1],
			}
		}
	}
}

// 提取错误码
func Code(err error) int32 {
	if err == nil {
		return 0
	}
	if e, ok := err.(*CodeError); ok {
		return e.Code
	}
	if e1, ok := err.(*TraceError); ok {
		if e2, ok := e1.err.(*CodeError); ok {
			return e2.Code
		}
	}
	s := strings.Split(err.Error(), ":")
	if len(s) == 0 {
		return CodeUnknown
	}
	i, err := strconv.Atoi(s[0])
	if err != nil {
		return CodeUnknown
	}
	return int32(i)
}

// 提取错误描述
func Msg(err error) string {
	if err == nil {
		return msg_OK
	}
	if e, ok := err.(*CodeError); ok {
		return e.Msg
	}
	if e1, ok := err.(*TraceError); ok {
		if e2, ok := e1.err.(*CodeError); ok {
			return e2.Msg
		}
	}
	s := strings.Split(err.Error(), ":")
	if len(s) < 2 {
		return err.Error()
	}
	return s[1]
}

type CodeError struct {
	Code   int32
	Msg    string
	Values map[string]interface{}
}

// 拷贝error并且附加stack
func (e *CodeError) Trace(values ...interface{}) *TraceError {
	var kvs map[string]interface{}
	if len(values)%2 != 0 {
		kvs = e.Values
	} else if len(values) == 0 {
		kvs = e.Values
	} else {
		kvs = make(map[string]interface{})
		for k, v := range e.Values {
			kvs[k] = v
		}
		for i := 0; i < len(values); i += 2 {
			j := i + 1
			if k, ok := values[i].(string); ok {
				kvs[k] = values[j]
			}
		}
	}
	return &TraceError{err: e, stack: debug.Stack(), values: kvs}
}

// 格式化错误字符串
func (e *CodeError) Error() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%d: %s", e.Code, e.Msg))
	if e.Values != nil {
		for k, v := range e.Values {
			sb.WriteString(fmt.Sprintf("\n%s: %s", k, v))
		}
	}
	return sb.String()
}

// 跟踪错误信息
type TraceError struct {
	err    error
	stack  []byte
	values map[string]interface{}
}

func (e *TraceError) Error() string {
	sb := strings.Builder{}
	sb.WriteString(e.err.Error())
	if e.values != nil {
		for k, v := range e.values {
			sb.WriteString(fmt.Sprintf("\n%s: %s", k, v))
		}
	}
	if e.stack != nil {
		sb.WriteString("\nstacktrace:\n")
		sb.Write(e.stack)
	}
	return sb.String()
}

func (e *TraceError) Unwrap() error {
	return e.err
}

func (e *TraceError) Is(err error) bool {
	if e.err == nil {
		return false
	}
	return e.err == err
}
