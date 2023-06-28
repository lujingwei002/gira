package codes

import (
	"bytes"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
)

// 错误码常量
const (
	CodeOk                             = 0
	CodeMalformed                      = -1 // error格式错误,不是"code:msg"格式
	CodeUnknown                        = -2
	CodeNullPointer                    = -2
	CodeNullObject                     = -7
	CodeInvalidSdkToken                = -18
	CodeInvalidArgs                    = -19
	CodeInvalidJwt                     = -20
	CodeJwtExpire                      = -21
	CodeToDo                           = -22
	CodeInvalidPassword                = -23
	CodePayOrderStatusInvalid          = -24
	CodeAccountPlatformNotSupport      = -25
	CodePayOrderAccountPlatformInvalid = -26
	CodePayOrderAmountInvalid          = -27
)
const (
	msg_OK                             = "成功"
	msg_Malformed                      = "malformed error"
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

func ThrowErrMalformed(values ...interface{}) *ErrorCode {
	return Throw(CodeMalformed, msg_Malformed, values...)
}

func ThrowErrUnknown(values ...interface{}) *ErrorCode {
	return Throw(CodeUnknown, msg_Unknown, values...)
}

func ThrowErrNullObject(values ...interface{}) *ErrorCode {
	return Throw(CodeNullObject, msg_NullObject, values...)
}

func ThrowErrInvalidSdkToken(values ...interface{}) *ErrorCode {
	return Throw(CodeInvalidSdkToken, msg_InvalidSdkToken, values...)
}

func ThrowErrInvalidArgs(values ...interface{}) *ErrorCode {
	return Throw(CodeInvalidArgs, msg_InvalidArgs, values...)
}

func ThrowErrInvalidJwt(values ...interface{}) *ErrorCode {
	return Throw(CodeInvalidJwt, msg_InvalidJwt, values...)
}

func ThrowErrJwtExpire(values ...interface{}) *ErrorCode {
	return Throw(CodeJwtExpire, msg_JwtExpire, values...)
}

func ThrowErrTodo(values ...interface{}) *ErrorCode {
	return Throw(CodeToDo, msg_ToDo, values...)
}

func ThrowErrInvalidPassword(values ...interface{}) *ErrorCode {
	return Throw(CodeInvalidPassword, msg_InvalidPassword, values...)
}

func ThrowErrPayOrderStatusInvalid(values ...interface{}) *ErrorCode {
	return Throw(CodePayOrderStatusInvalid, msg_PayOrderStatusInvalid, values...)
}

func ThrowErrAccountPlatformNotSupport(values ...interface{}) *ErrorCode {
	return Throw(CodeAccountPlatformNotSupport, msg_AccountPlatformNotSupport, values...)
}

func ThrowErrPayOrderAccountPlatformInvalid(values ...interface{}) *ErrorCode {
	return Throw(CodePayOrderAccountPlatformInvalid, msg_PayOrderAccountPlatformInvalid, values...)
}

func ThrowErrPayOrderAmountInvalid(values ...interface{}) *ErrorCode {
	return Throw(CodePayOrderAmountInvalid, msg_PayOrderAmountInvalid, values...)
}

// 根据code, msg, values创建error_code
func Throw(code int32, msg string, values ...interface{}) *ErrorCode {
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
	os.Stderr.WriteString(fmt.Sprintf("%s\n", kvs))
	debug.PrintStack()
	return &ErrorCode{
		Code:   code,
		Msg:    msg,
		Values: kvs,
	}
}

// 根据code, msg, values创建error_code
func New(code int32, msg string, values ...interface{}) *ErrorCode {
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
				Code: CodeMalformed,
				Msg:  msg,
			}
		} else if v, err := strconv.Atoi(arr[0]); err != nil {
			return &ErrorCode{
				Code: CodeMalformed,
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
func Code(err error) int32 {
	if err == nil {
		return 0
	}
	if e, ok := err.(*ErrorCode); ok {
		return e.Code
	}
	s := strings.Split(err.Error(), ":")
	if len(s) == 0 {
		return CodeMalformed
	}
	i, err := strconv.Atoi(s[0])
	if err != nil {
		return CodeMalformed
	}
	return int32(i)
}

// 提取错误描述
func Msg(err error) string {
	if err == nil {
		return msg_OK
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
func Stack(err error) string {
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
