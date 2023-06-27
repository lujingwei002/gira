package gira

import (
	"errors"
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
	ErrNullPointer                        = errors.New("null pointer")
	ErrResourceManagerNotImplement        = errors.New("resource manager not implement")
	ErrConfigHandlerNotImplement          = errors.New("config handler not implement")
	ErrResourceLoaderNotImplement         = errors.New("resource loader not implement")
	ErrResourceHandlerNotImplement        = errors.New("resource handler not implement")
	ErrHttpHandlerNotImplement            = errors.New("http handler not implement")
	ErrRegisterServerFail                 = errors.New("register server fail")
	ErrInvalidPeer                        = errors.New("invalid peer")
	ErrDataUpsertFail                     = errors.New("upsert fail")
	ErrDataExist                          = errors.New("data exist")
	ErrDataNotExist                       = errors.New("data not exist")
	ErrDataNotFound                       = errors.New("data not found")
	ErrDataInsertFail                     = errors.New("data insert fail")
	ErrDataDeleteFail                     = errors.New("data delete fail")
	ErrSdkComponentNotImplement           = errors.New("sdk commponent not implement")
	ErrGateHandlerNotImplement            = errors.New("gate handler not implement")
	ErrPeerHandlerNotImplement            = errors.New("peer handler not implement")
	ErrGrpcServerNotImplement             = errors.New("grpc handler not implement")
	ErrHallHandlerNotImplement            = errors.New("hall handler not implement")
	ErrSprotoHandlerNotImplement          = errors.New("sproto handler not implement")
	ErrSprotoReqIdConflict                = errors.New("sproto req id conflict")
	ErrSprotoReqTimeout                   = errors.New("sproto req id timeout")
	ErrSprotoResponseConversion           = errors.New("sproto response type conversion")
	ErrReadOnClosedClient                 = errors.New("read on closed client")
	ErrPeerNotFound                       = errors.New("peer not found")
	ErrUserInstead                        = errors.New("账号在其他地方登录")
	ErrUserLocked                         = errors.New("账号在其他地方被锁定")
	ErrGrpcClientPoolNil                  = errors.New("grpc pool无法申请client")
	ErrBrokenChannel                      = errors.New("管道已关闭，不能再写数据")
	ErrSessionClosed                      = errors.New("会话已经关闭")
	ErrUpstreamUnavailable                = errors.New("上游服务不可用")
	ErrUpstreamUnreachable                = errors.New("上游服务不可达")
	ErrServiceUnavailable                 = errors.New("服务不可用")
	ErrSprotoPushConversion               = errors.New("sproto push type conversion")
	ErrProjectFileNotFound                = errors.New("gira.yaml文件找不到")
	ErrInvalidPassword                    = errors.New("密码错误")
	ErrGenNotChange                       = errors.New("gen源文件没变化")
	ErrNoSession                          = errors.New("玩家不在线")
	ErrSprotoWaitPushTimeout              = errors.New("等待push超时")
	ErrSprotoWaitPushConflict             = errors.New("等待push冲突 ，只可以有一个")
	ErrTodo                               = &TodoError{}
	ErrServerDown                         = errors.New("服务器关闭")
	ErrGrpcServerNotOpen                  = errors.New("grpc模块未开启")
	ErrAdminClientNotImplement            = errors.New("admin client 接口末实现")
	ErrRegistryNOtImplement               = errors.New("注册表功能未实现")
	ErrInvalidMemberId                    = errors.New("member id非法")
	ErrBehaviorNotInit                    = errors.New("behavior driver not init")
	ErrDbNotSupport                       = errors.New("数据库类型不支持")
	ErrDbNotConfig                        = errors.New("数据库配置异常")
	ErrInvalidService                     = errors.New("service格式非法")
	ErrServiceNotFound                    = errors.New("查找不到service")
	ErrServiceLocked                      = errors.New("注册service失败")
	ErrAlreadyDestory                     = errors.New("已经销毁")
	ErrServiceContainerNotImplement       = errors.New("service container 未实现")
	ErrServiceNotImplement                = errors.New("service接口未实现")
	ErrUserNotFound                       = errors.New("用户不在线")
	ErrActorCallTimeOut                   = errors.New("actor call timeout")
	ErrEtcdConfigNotFound                 = errors.New("etcd未配置")
	ErrGrpcConfigNotFound                 = errors.New("grpc未配置")
	ErrInterrupt                          = errors.New("interrupt")
	ErrServerRouterMetaNotFound           = errors.New("server router incoming context not found")
	ErrServerRouterKeyNotFound            = errors.New("server router-key not found")
	ErrServerRouterHandlerNotRegist       = errors.New("server router handler not regist")
	ErrServerRouterHandlerNotImplement    = errors.New("cata server handler not implement")
	ErrServiceAlreadyStopped              = errors.New("service已经停止")
	ErrServiceAlreadyStarted              = errors.New("service已经启动")
	ErrDbClientComponentNotImplement      = errors.New("db client component not implement")
	ErrServerNotFound                     = errors.New("server not found")
	ErrSdkPayOrderCheckMethodNotImplement = errors.New("sdk pay order check 方法未实现")
	ErrSdkPayOrderCheckArgsInvalid        = errors.New("sdk pay order check 参数错误")
	errAccountPlatformNotSupport          = errors.New("account platform 不支持")
	errPayOrderStatusInvalid              = errors.New("pay order 状态错误")
	errPayOrderAmountInvalid              = errors.New("pay order 价格错误")
	errPayOrderAccountPlatformInvalid     = errors.New("pay order 账号平台错误")
	ErrCronNotImplement                   = errors.New("cron not implement")
)

func NewPeerAlreadyRegistError(key string) error {
	return fmt.Errorf("peer already regist: %s", key)
}

type Error struct {
	Msg    string
	Stack  []byte
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
	if e.Stack != nil {
		sb.WriteString("\ntraceback:\n")
		sb.Write(e.Stack)
	}
	return sb.String()
}

func Errorw(msg string, values ...interface{}) *Error {
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
