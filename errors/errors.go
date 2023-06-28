package errors

import (
	"fmt"
	"runtime/debug"
	"strings"
)

type SyntaxError struct {
	msg      string
	filePath string
	lines    []string
}

func NewSyntaxError(msg string, filePath string, content string) *SyntaxError {
	lines := strings.Split(content, "\n")
	return &SyntaxError{
		msg:      msg,
		filePath: filePath,
		lines:    lines,
	}
}

func (e *SyntaxError) Error() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s\n", e.msg))
	sb.WriteString(fmt.Sprintf("%s\n", e.filePath))
	for k, v := range e.lines {
		sb.WriteString(fmt.Sprintf("%d: %s\n", k+1, v))
	}
	return sb.String()
}

type TodoError struct {
	stack []byte
}

func (e *TodoError) Error() string {
	sb := strings.Builder{}
	return sb.String()
}

func (e *TodoError) Trace() error {
	e.stack = debug.Stack()
	return e
}

var (
	ErrTodo                               = &TodoError{}
	ErrNullObject                         = New("null object")
	ErrNullPointer                        = New("null pointer")
	ErrInvalidArgs                        = New("invalid args")
	ErrResourceManagerNotImplement        = New("resource manager not implement")
	ErrConfigHandlerNotImplement          = New("config handler not implement")
	ErrResourceLoaderNotImplement         = New("resource loader not implement")
	ErrResourceHandlerNotImplement        = New("resource handler not implement")
	ErrHttpHandlerNotImplement            = New("http handler not implement")
	ErrSdkComponentNotImplement           = New("sdk commponent not implement")
	ErrGateHandlerNotImplement            = New("gate handler not implement")
	ErrPeerHandlerNotImplement            = New("peer handler not implement")
	ErrGrpcServerNotImplement             = New("grpc handler not implement")
	ErrHallHandlerNotImplement            = New("hall handler not implement")
	ErrSprotoHandlerNotImplement          = New("sproto handler not implement")
	ErrDbClientComponentNotImplement      = New("db client component not implement")
	ErrServerRouterHandlerNotImplement    = New("cata server handler not implement")
	ErrCronNotImplement                   = New("cron not implement")
	ErrAdminClientNotImplement            = New("admin client 接口末实现")
	ErrRegistryNOtImplement               = New("注册表功能未实现")
	ErrServiceContainerNotImplement       = New("service container 未实现")
	ErrSdkPayOrderCheckMethodNotImplement = New("sdk pay order check 方法未实现")
	ErrServiceNotImplement                = New("service接口未实现")
	ErrDataNotExist                       = New("data not exist")
	ErrDataNotFound                       = New("data not found")
	ErrDataInsertFail                     = New("data insert fail")
	ErrDataDeleteFail                     = New("data delete fail")
	ErrPeerNotFound                       = New("peer not found")
	ErrUserInstead                        = New("账号在其他地方登录")
	ErrUserLocked                         = New("账号在其他地方被锁定")
	ErrGrpcClientPoolNil                  = New("grpc pool无法申请client")
	ErrBrokenChannel                      = New("管道已关闭，不能再写数据")
	ErrSessionClosed                      = New("会话已经关闭")
	ErrUpstreamUnavailable                = New("上游服务不可用")
	ErrUpstreamUnreachable                = New("上游服务不可达")
	ErrServiceUnavailable                 = New("服务不可用")
	ErrProjectFileNotFound                = New("gira.yaml文件找不到")
	ErrInvalidPassword                    = New("密码错误")
	ErrGenNotChange                       = New("gen源文件没变化")
	ErrNoSession                          = New("玩家不在线")
	ErrInvalidMemberId                    = New("member id非法")
	ErrBehaviorNotInit                    = New("behavior driver not init")
	ErrDbNotSupport                       = New("数据库类型不支持")
	ErrInvalidService                     = New("service格式非法")
	ErrServiceNotFound                    = New("查找不到service")
	ErrServiceLocked                      = New("注册service失败")
	ErrUserNotFound                       = New("用户不在线")
	ErrActorCallTimeOut                   = New("actor call timeout")
	ErrInterrupt                          = New("interrupt")
	ErrServerRouterMetaNotFound           = New("server router incoming context not found")
	ErrServerRouterKeyNotFound            = New("server router-key not found")
	ErrServerRouterHandlerNotRegist       = New("server router handler not regist")
	ErrServerNotFound                     = New("server not found")
	ErrInvalidJwt                         = New("invalid jwt")
	ErrJwtExpire                          = New("jwt expire")
	ErrInvalidSdkToken                    = New("invalid sdk token")
)

type Error struct {
	Msg    string
	Values map[string]interface{}
}

// 格式化错误字符串
func (e *Error) Error() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s\n", e.Msg))
	if e.Values != nil {
		for k, v := range e.Values {
			sb.WriteString(fmt.Sprintf("\n%s: %s", k, v))
		}
	}
	return sb.String()
}

func New(msg string, values ...interface{}) *Error {
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
	return &Error{
		Msg:    msg,
		Values: kvs,
	}
}

func Throw(msg string, values ...interface{}) *Error {
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
	return &Error{
		Msg:    msg,
		Values: kvs,
	}
}
