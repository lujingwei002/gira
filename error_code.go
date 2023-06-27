package gira

import (
	"bytes"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
)

// 错误码常量
const (
	E_OK                = 0
	E_MALFORMED         = -1 // error格式错误,不是"code:msg"格式
	E_UNKNOWN           = -2
	E_NULL_POINTER      = -2
	E_NULL_OBJECT       = -7
	E_INVALID_SDK_TOKEN = -18
	E_INVALID_ARGS      = -19
	E_INVALID_JWT       = -20
	E_JWT_EXPIRE        = -21
)
const (
	E_MSG_OK                = "成功"
	E_MSG_MALFORMED         = "malformed error"
	E_MSG_UNKNOWN           = "未知错误"
	E_MSG_NULL_OBJECT       = "null object"
	E_MSG_INVALID_SDK_TOKEN = "无效的sdk token"
	E_MSG_INVALID_ARGS      = "无效参数"
	E_MSG_INVALID_JWT       = "无效的token"
	E_MSG_JWT_EXPIRE        = "token已过期"
)

var (
	ErrMalformed       = NewErrorCode(E_MALFORMED, E_MSG_MALFORMED)
	ErrUnknown         = NewErrorCode(E_UNKNOWN, E_MSG_UNKNOWN)
	ErrNullObject      = NewErrorCode(E_NULL_OBJECT, E_MSG_NULL_OBJECT)
	ErrInvalidSdkToken = NewErrorCode(E_INVALID_SDK_TOKEN, E_MSG_INVALID_SDK_TOKEN)
	ErrInvalidArgs     = NewErrorCode(E_INVALID_ARGS, E_MSG_INVALID_ARGS)
	ErrInvalidJwt      = NewErrorCode(E_INVALID_JWT, E_MSG_INVALID_JWT)
	ErrJwtExpire       = NewErrorCode(E_JWT_EXPIRE, E_MSG_JWT_EXPIRE)
)

// 根据code, msg创建error_code
func NewErrorCode(code int32, msg string) *ErrorCode {
	return &ErrorCode{
		Code: code,
		Msg:  msg,
	}
}

// 根据code, msg, values创建error_code
func NewErrorCodew(code int32, msg string, values ...interface{}) *ErrorCode {
	var kvs map[string]interface{}
	for i := 0; i < len(values); i += 2 {
		j := i + 1
		if j < len(values) {
			if kvs == nil {
				kvs = make(map[string]interface{})
			}
			if k, ok := values[i].(string); ok {
				kvs[k] = values[j]
			}
		}
	}
	return &ErrorCode{
		Code:   code,
		Msg:    msg,
		Values: kvs,
	}
}

// 转换或者创建error
func NewOrCastError(err error) *ErrorCode {
	if err == nil {
		return nil
	}
	if e, ok := err.(*ErrorCode); ok {
		return e
	} else {
		msg := err.Error()
		arr := strings.Split(msg, ":")
		if len(arr) < 2 {
			return &ErrorCode{
				Code: E_MALFORMED,
				Msg:  msg,
			}
		} else if v, err := strconv.Atoi(arr[0]); err != nil {
			return &ErrorCode{
				Code: E_MALFORMED,
				Msg:  msg,
			}
		} else {
			return &ErrorCode{
				Code: int32(v),
				Msg:  arr[1],
			}
		}
	}
}

// 提取错误码
func ErrCode(err error) int32 {
	if err == nil {
		return 0
	}
	if e, ok := err.(*ErrorCode); ok {
		return e.Code
	}
	s := strings.Split(err.Error(), ":")
	if len(s) == 0 {
		return E_MALFORMED
	}
	i, err := strconv.Atoi(s[0])
	if err != nil {
		return E_MALFORMED
	}
	return int32(i)
}

// 提取错误描述
func ErrMsg(err error) string {
	if err == nil {
		return E_MSG_OK
	}
	if e, ok := err.(*ErrorCode); ok {
		return e.Msg
	}
	s := strings.Split(err.Error(), ":")
	if len(s) < 2 {
		return err.Error()
	}
	return s[1]
}

// 提取错误栈
func ErrStack(err error) string {
	if e, ok := err.(*ErrorCode); ok {
		return string(e.Stack)
	}
	return ""
}

type ErrorCode struct {
	Code   int32
	Msg    string
	Stack  []byte
	Values map[string]interface{}
}

// 拷贝error并且附加stack
func (e *ErrorCode) Trace() *ErrorCode {
	return &ErrorCode{Code: e.Code, Msg: e.Msg, Stack: debug.Stack(), Values: e.Values}
}

// 格式化错误字符串
func (e *ErrorCode) Error() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%d: %s", e.Code, e.Msg))
	if e.Values != nil {
		for k, v := range e.Values {
			sb.WriteString(fmt.Sprintf("\n%s: %s", k, v))
		}
	}
	if e.Stack != nil {
		sb.WriteString("\ntraceback:\n")
		sb.Write(e.Stack)
	}
	return sb.String()
}

// 附加value到error中
func (err *ErrorCode) WithValues(values ...interface{}) *ErrorCode {
	for i := 0; i < len(values); i += 2 {
		j := i + 1
		if j < len(values) {
			if err.Values == nil {
				err.Values = make(map[string]interface{})
			}
			if k, ok := values[i].(string); ok {
				err.Values[k] = values[j]
			}
		}
	}
	return err
}

// 附加stack到error中
func (err *ErrorCode) WithTrace(values ...interface{}) *ErrorCode {
	err.Stack = debug.Stack()
	return err
}

// 附加错误描述到error中
func (err *ErrorCode) WithMsg(msg string) *ErrorCode {
	err.Msg = msg
	return err
}

func (err *ErrorCode) WithFile(file string) *ErrorCode {
	if err.Values == nil {
		err.Values = make(map[string]interface{})
	}
	err.Values["file"] = file
	return err
}

func (err *ErrorCode) WithLines(content []byte) *ErrorCode {
	sb := strings.Builder{}
	arr := bytes.Split(content, []byte("\n"))
	for index, line := range arr {
		sb.WriteString(fmt.Sprintf("%d: %s\n", index+1, line))
	}
	if err.Values == nil {
		err.Values = make(map[string]interface{})
	}
	err.Values["lines"] = sb.String()
	return err
}
