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
	E_OK                                       = 0
	E_MALFORMED                                = -1 // error格式错误,不是"code:msg"格式
	E_UNKNOWN                                  = -2
	E_NULL_POINTER                             = -2
	E_RESOURCE_MANAGER_NOT_IMPLEMENT           = -3
	E_CONFIG_HANDLER_NOT_IMPLEMENT             = -4
	E_RESOURCE_LOADER_NOT_IMPLEMENT            = -5
	E_RESOURCE_HANDLER_NOT_IMPLEMENT           = -6
	E_NULL_OBJECT                              = -7
	E_HTTP_HANDLER_NOT_IMPLEMENT               = -8
	E_REGISTER_SERVER_FAIL                     = -9
	E_INVALID_PEER                             = -10
	E_DATA_UPSERT_FAIL                         = -11
	E_DATA_EXIST                               = -12
	E_DATA_NOT_EXIST                           = -13
	E_DATA_NOT_FOUND                           = -14
	E_DATA_INSERT_FAIL                         = -15
	E_DATA_DELETE_FAIL                         = -16
	E_SDK_COMPONENT_NOT_IMPLEMENT              = -17
	E_INVALID_SDK_TOKEN                        = -18
	E_INVALID_ARGS                             = -19
	E_INVALID_JWT                              = -20
	E_JWT_EXPIRE                               = -21
	E_GATE_HANDLER_NOT_IMPLEMENT               = -22
	E_PEER_HANDLER_NOT_IMPLEMENT               = -23
	E_GRPC_SERVER_NOT_IMPLEMENT                = -24
	E_HALL_HANDLER_NOT_IMPLEMENT               = -25
	E_SPROTO_HANDLER_NOT_IMPLEMENT             = -26
	E_SPROTO_REQ_ID_CONFLICT                   = -27
	E_SPROTO_REQ_TIMEOUT                       = -28
	E_SPROTO_RESPONSE_CONVERSION               = -29
	E_READ_ON_CLOSED_CLIENT                    = -30
	E_PEER_NOT_FOUND                           = -31
	E_USER_INSTEAD                             = -32
	E_USER_LOCKED                              = -33
	E_GRPC_CLIENT_POOL_NIL                     = -34
	E_BROKEN_CHANNEL                           = -35
	E_SESSION_CLOSED                           = -36
	E_UPSTREAM_UNAVAILABLE                     = -37
	E_UPSTREAM_UNREACHABLE                     = -38
	E_SERVICE_UNAVAILABLE                      = -39
	E_SPROTO_PUSH_CONVERSION                   = -40
	E_PROJECT_FILE_NOT_FOUND                   = -41
	E_INVALID_PASSWORD                         = -42
	E_GEN_NOT_CHANGE                           = -43
	E_NO_SESSION                               = -44
	E_SPROTO_WAIT_PUSH_TIMEOUT                 = -45
	E_SPROTO_WAIT_PUSH_CONFLICT                = -46
	E_TODO                                     = -47
	E_SERVER_DOWN                              = -48
	E_GRPC_SERVER_NOT_OPEN                     = -49
	E_ADMIN_CLIENT_NOT_IMPLEMENT               = -50
	E_REGISTRY_NOT_IMPLEMENT                   = -51
	E_INVALID_MEMBER_ID                        = -52
	E_BEHAVIOR_DRIVER_NOT_INIT                 = -53
	E_DB_NOT_SUPPORT                           = -54
	E_DB_NOT_CONFIG                            = -55
	E_INVALID_SERVICE                          = -56
	E_SERVICE_NOT_FOUND                        = -57
	E_SERVICE_LOCKED                           = -58
	E_ALREADY_DESTORY                          = -59
	E_SERVICE_CONTAINER_NOT_IMPLEMENT          = -59
	E_SERVICE_NOT_IMPLEMENT                    = -59
	E_USER_NOT_FOUND                           = -60
	E_ACTOR_CALL_TIME_OUT                      = -61
	E_ETCD_CONFIG_NOT_FOUND                    = -62
	E_GRPC_CONFIG_NOT_FOUND                    = -63
	E_INTERRUPT                                = -64
	E_CATALOG_SERVER_META_NOT_FOUND            = -65
	E_CATALOG_SERVER_KEY_NOT_FOUND             = -66
	E_CATALOG_SERVER_HANDLER_NOT_REGIST        = -67
	E_CATALOG_SERVER_HANDLER_NOT_IMPLEMENT     = -68
	E_SERVICE_ALREADY_STOPPED                  = -69
	E_SERVICE_ALREADY_STARTED                  = -70
	E_DB_CLIENT_COMPONENT_NOT_IMPLEMENT        = -71
	E_SERVER_NOT_FOUND                         = -72
	E_SDK_PAY_ORDER_CHECK_METHOD_NOT_IMPLEMENT = -73
	E_SDK_PAY_ORDER_CHECK_ARGS_INVALID         = -74
	E_ACCOUNT_PLATFORM_NOT_SUPPORT             = -75
	E_PAY_ORDER_STATUS_INVALID                 = -76
	E_PAY_ORDER_AMOUNT_INVALID                 = -77
	E_PAY_ORDER_ACCOUNT_PLATFORM_INVALID       = -78
	E_CAST                                     = -79
	E_PROTO_REQUEST_CAST                       = -80
	E_PROTO_RESPONSE_CAST                      = -81
	E_PROTO_PUSH_CAST                          = -82
	E_PROTO_RESPONSE_NEW                       = -83
	E_CRON_NOT_IMPLEMENT                       = -84
	E_CONFIG_ENV_INVALID_SYNTAX                = -85
	E_INVALID_SYNTAX                           = -86
)
const (
	E_MSG_OK                                       = "成功"
	E_MSG_MALFORMED                                = "malformed error"
	E_MSG_UNKNOWN                                  = "未知错误"
	E_MSG_NULL_POINTER                             = "空指针"
	E_MSG_RESOURCE_MANAGER_NOT_IMPLEMENT           = "resource manager not implement"
	E_MSG_CONFIG_HANDLER_NOT_IMPLEMENT             = "config handler not implement"
	E_MSG_RESOURCE_LOADER_NOT_IMPLEMENT            = "resource loader not implement"
	E_MSG_RESOURCE_HANDLER_NOT_IMPLEMENT           = "resource handler not implement"
	E_MSG_NULL_OBJECT                              = "null object"
	E_MSG_HTTP_HANDLER_NOT_IMPLEMENT               = "http handler not implement"
	E_MSG_REGISTER_SERVER_FAIL                     = "register server fail"
	E_MSG_INVALID_PEER                             = "invalid peer"
	E_MSG_DATA_UPSERT_FAIL                         = "upsert fail"
	E_MSG_DATA_EXIST                               = "data exist"
	E_MSG_DATA_NOT_EXIST                           = "data not exist"
	E_MSG_DATA_NOT_FOUND                           = "data not found"
	E_MSG_DATA_INSERT_FAIL                         = "data insert fail"
	E_MSG_DATA_DELETE_FAIL                         = "data delete fail"
	E_MSG_SDK_COMPONENT_NOT_IMPLEMENT              = "sdk commponent not implement"
	E_MSG_INVALID_SDK_TOKEN                        = "无效的sdk token"
	E_MSG_INVALID_ARGS                             = "无效参数"
	E_MSG_INVALID_JWT                              = "无效的token"
	E_MSG_JWT_EXPIRE                               = "token已过期"
	E_MSG_GATE_HANDLER_NOT_IMPLEMENT               = "gate handler not implement"
	E_MSG_PEER_HANDLER_NOT_IMPLEMENT               = "peer handler not implement"
	E_MSG_GRPC_SERVER_NOT_IMPLEMENT                = "grpc handler not implement"
	E_MSG_HALL_HANDLER_NOT_IMPLEMENT               = "hall handler not implement"
	E_MSG_SPROTO_HANDLER_NOT_IMPLEMENT             = "sproto handler not implement"
	E_MSG_SPROTO_REQ_ID_CONFLICT                   = "sproto req id conflict"
	E_MSG_SPROTO_REQ_TIMEOUT                       = "sproto req id timeout"
	E_MSG_SPROTO_RESPONSE_CONVERSION               = "sproto response type conversion"
	E_MSG_READ_ON_CLOSED_CLIENT                    = "read on closed client"
	E_MSG_PEER_NOT_FOUND                           = "peer not found"
	E_MSG_USER_INSTEAD                             = "账号在其他地方登录"
	E_MSG_USER_LOCKED                              = "账号在其他地方被锁定"
	E_MSG_GRPC_CLIENT_POOL_NIL                     = "grpc pool无法申请client"
	E_MSG_BROKEN_CHANNEL                           = "管道已关闭，不能再写数据"
	E_MSG_SESSION_CLOSED                           = "会话已经关闭"
	E_MSG_UPSTREAM_UNAVAILABLE                     = "上游服务不可用"
	E_MSG_UPSTREAM_UNREACHABLE                     = "上游服务不可达"
	E_MSG_SERVICE_UNAVAILABLE                      = "服务不可用"
	E_MSG_SPROTO_PUSH_CONVERSION                   = "sproto push type conversion"
	E_MSG_PROJECT_FILE_NOT_FOUND                   = "gira.yaml文件找不到"
	E_MSG_INVALID_PASSWORD                         = "密码错误"
	E_MSG_GEN_NOT_CHANGE                           = "gen源文件没变化"
	E_MSG_NO_SESSION                               = "玩家不在线"
	E_MSG_SPROTO_WAIT_PUSH_TIMEOUT                 = "等待push超时"
	E_MSG_SPROTO_WAIT_PUSH_CONFLICT                = "等待push冲突 ，只可以有一个"
	E_MSG_TODO                                     = "TODO"
	E_MSG_SERVER_DOWN                              = "服务器关闭"
	E_MSG_GRPC_SERVER_NOT_OPEN                     = "grpc模块未开启"
	E_MSG_ADMIN_CLIENT_NOT_IMPLEMENT               = "admin client 接口末实现"
	E_MSG_REGISTRY_NOT_IMPLEMENT                   = "注册表功能未实现"
	E_MSG_INVALID_MEMBER_ID                        = "member id非法"
	E_MSG_BEHAVIOR_DRIVER_NOT_INIT                 = "behavior driver not init"
	E_MSG_DB_NOT_SUPPORT                           = "数据库类型不支持"
	E_MSG_DB_NOT_CONFIG                            = "数据库配置异常"
	E_MSG_INVALID_SERVICE                          = "service格式非法"
	E_MSG_SERVICE_NOT_FOUND                        = "查找不到service"
	E_MSG_SERVICE_LOCKED                           = "注册service失败"
	E_MSG_ALREADY_DESTORY                          = "已经销毁"
	E_MSG_SERVICE_ONTAINER_NOT_IMPLEMENT           = "service container 未实现"
	E_MSG_SERVICE_NOT_IMPLEMENT                    = "service接口未实现"
	E_MSG_USER_NOT_FOUND                           = "用户不在线"
	E_MSG_ACTOR_CALL_TIME_OUT                      = "actor call timeout"
	E_MSG_ETCD_CONFIG_NOT_FOUND                    = "etcd未配置"
	E_MSG_GRPC_CONFIG_NOT_FOUND                    = "grpc未配置"
	E_MSGS_INTERRUPT                               = "interrupt"
	E_MSG_CATALOG_SERVER_META_NOT_FOUND            = "catalog server incoming context not found"
	E_MSG_CATALOG_SERVER_KEY_NOT_FOUND             = "catalog server catalog-key not found"
	E_MSG_CATALOG_SERVER_HANDLER_NOT_REGIST        = "catalog server handler not regist"
	E_MSG_CATALOG_SERVER_HANDLER_NOT_IMPLEMENT     = "cata server handler not implement"
	E_MSG_SERVICE_ALREADY_STOPPED                  = "service已经停止"
	E_MSG_SERVICE_ALREADY_STARTED                  = "service已经启动"
	E_MSG_DB_CLIENT_COMPONENT_NOT_IMPLEMENT        = "db client component not implement"
	E_MSG_SERVER_NOT_FOUND                         = "server not found"
	E_MSG_SDK_PAY_ORDER_CHECK_METHOD_NOT_IMPLEMENT = "sdk pay order check 方法未实现"
	E_MSG_SDK_PAY_ORDER_CHECK_ARGS_INVALID         = "sdk pay order check 参数错误"
	E_MSG_ACCOUNT_PLATFORM_NOT_SUPPORT             = "account platform 不支持"
	E_MSG_PAY_ORDER_STATUS_INVALID                 = "pay order 状态错误"
	E_MSG_PAY_ORDER_AMOUNT_INVALID                 = "pay order 价格错误"
	E_MSG_PAY_ORDER_ACCOUNT_PLATFORM_INVALID       = "pay order 账号平台错误"
	E_MSG_CAST                                     = "cast"
	E_MSG_PROTO_REQUEST_CAST                       = "proto request cast"
	E_MSG_PROTO_RESPONSE_CAST                      = "proto response cast"
	E_MSG_PROTO_PUSH_CAST                          = "proto push cast"
	E_MSG_PROTO_RESPONSE_NEW                       = "proto response new"
	E_MSG_CRON_NOT_IMPLEMENT                       = "cron not implement "
	E_MSG_CONFIG_ENV_INVALID_SYNTAX                = "config env file invalid syntax"
	E_MSG_INVALID_SYNTAX                           = "invalid syntax"
)

var (
	ErrNullPonter                         = NewError(E_RESOURCE_MANAGER_NOT_IMPLEMENT, E_MSG_NULL_POINTER)
	ErrMalformed                          = NewError(E_MALFORMED, E_MSG_MALFORMED)
	ErrUnknown                            = NewError(E_UNKNOWN, E_MSG_UNKNOWN)
	ErrResourceManagerNotImplement        = NewError(E_CONFIG_HANDLER_NOT_IMPLEMENT, E_MSG_RESOURCE_MANAGER_NOT_IMPLEMENT)
	ErrConfigHandlerNotImplement          = NewError(E_CONFIG_HANDLER_NOT_IMPLEMENT, E_MSG_CONFIG_HANDLER_NOT_IMPLEMENT)
	ErrResourceLoaderNotImplement         = NewError(E_RESOURCE_LOADER_NOT_IMPLEMENT, E_MSG_RESOURCE_LOADER_NOT_IMPLEMENT)
	ErrResourceHandlerNotImplement        = NewError(E_RESOURCE_HANDLER_NOT_IMPLEMENT, E_MSG_RESOURCE_HANDLER_NOT_IMPLEMENT)
	ErrNullObject                         = NewError(E_NULL_OBJECT, E_MSG_NULL_OBJECT)
	ErrHttpHandlerNotImplement            = NewError(E_HTTP_HANDLER_NOT_IMPLEMENT, E_MSG_HTTP_HANDLER_NOT_IMPLEMENT)
	ErrRegisterServerFail                 = NewError(E_REGISTER_SERVER_FAIL, E_MSG_REGISTER_SERVER_FAIL)
	ErrInvalidPeer                        = NewError(E_INVALID_PEER, E_MSG_INVALID_PEER)
	ErrDataUpsertFail                     = NewError(E_DATA_UPSERT_FAIL, E_MSG_DATA_UPSERT_FAIL)
	ErrDataExist                          = NewError(E_DATA_EXIST, E_MSG_DATA_EXIST)
	ErrDataNotExist                       = NewError(E_DATA_NOT_EXIST, E_MSG_DATA_NOT_EXIST)
	ErrDataNotFound                       = NewError(E_DATA_NOT_FOUND, E_MSG_DATA_NOT_FOUND)
	ErrDataInsertFail                     = NewError(E_DATA_INSERT_FAIL, E_MSG_DATA_INSERT_FAIL)
	ErrDataDeleteFail                     = NewError(E_DATA_DELETE_FAIL, E_MSG_DATA_DELETE_FAIL)
	ErrSdkComponentNotImplement           = NewError(E_SDK_COMPONENT_NOT_IMPLEMENT, E_MSG_SDK_COMPONENT_NOT_IMPLEMENT)
	ErrInvalidSdkToken                    = NewError(E_INVALID_SDK_TOKEN, E_MSG_INVALID_SDK_TOKEN)
	ErrInvalidArgs                        = NewError(E_INVALID_ARGS, E_MSG_INVALID_ARGS)
	ErrInvalidJwt                         = NewError(E_INVALID_JWT, E_MSG_INVALID_JWT)
	ErrJwtExpire                          = NewError(E_JWT_EXPIRE, E_MSG_JWT_EXPIRE)
	ErrGateHandlerNotImplement            = NewError(E_GATE_HANDLER_NOT_IMPLEMENT, E_MSG_GATE_HANDLER_NOT_IMPLEMENT)
	ErrPeerHandlerNotImplement            = NewError(E_PEER_HANDLER_NOT_IMPLEMENT, E_MSG_PEER_HANDLER_NOT_IMPLEMENT)
	ErrGrpcServerNotImplement             = NewError(E_GRPC_SERVER_NOT_IMPLEMENT, E_MSG_GRPC_SERVER_NOT_IMPLEMENT)
	ErrHallHandlerNotImplement            = NewError(E_HALL_HANDLER_NOT_IMPLEMENT, E_MSG_HALL_HANDLER_NOT_IMPLEMENT)
	ErrSprotoHandlerNotImplement          = NewError(E_SPROTO_HANDLER_NOT_IMPLEMENT, E_MSG_SPROTO_HANDLER_NOT_IMPLEMENT)
	ErrSprotoReqIdConflict                = NewError(E_SPROTO_REQ_ID_CONFLICT, E_MSG_SPROTO_REQ_ID_CONFLICT)
	ErrSprotoReqTimeout                   = NewError(E_SPROTO_REQ_TIMEOUT, E_MSG_SPROTO_REQ_TIMEOUT)
	ErrSprotoResponseConversion           = NewError(E_SPROTO_RESPONSE_CONVERSION, E_MSG_SPROTO_RESPONSE_CONVERSION)
	ErrReadOnClosedClient                 = NewError(E_READ_ON_CLOSED_CLIENT, E_MSG_READ_ON_CLOSED_CLIENT)
	ErrPeerNotFound                       = NewError(E_PEER_NOT_FOUND, E_MSG_PEER_NOT_FOUND)
	ErrUserInstead                        = NewError(E_USER_INSTEAD, E_MSG_USER_INSTEAD)
	ErrUserLocked                         = NewError(E_USER_LOCKED, E_MSG_USER_LOCKED)
	ErrGrpcClientPoolNil                  = NewError(E_GRPC_CLIENT_POOL_NIL, E_MSG_GRPC_CLIENT_POOL_NIL)
	ErrBrokenChannel                      = NewError(E_BROKEN_CHANNEL, E_MSG_BROKEN_CHANNEL)
	ErrSessionClosed                      = NewError(E_SESSION_CLOSED, E_MSG_SESSION_CLOSED)
	ErrUpstreamUnavailable                = NewError(E_UPSTREAM_UNAVAILABLE, E_MSG_UPSTREAM_UNAVAILABLE)
	ErrUpstreamUnreachable                = NewError(E_UPSTREAM_UNREACHABLE, E_MSG_UPSTREAM_UNREACHABLE)
	ErrServiceUnavailable                 = NewError(E_SERVICE_UNAVAILABLE, E_MSG_SERVICE_UNAVAILABLE)
	ErrSprotoPushConversion               = NewError(E_SPROTO_PUSH_CONVERSION, E_MSG_SPROTO_PUSH_CONVERSION)
	ErrProjectFileNotFound                = NewError(E_PROJECT_FILE_NOT_FOUND, E_MSG_PROJECT_FILE_NOT_FOUND)
	ErrInvalidPassword                    = NewError(E_INVALID_PASSWORD, E_MSG_INVALID_PASSWORD)
	ErrGenNotChange                       = NewError(E_GEN_NOT_CHANGE, E_MSG_GEN_NOT_CHANGE)
	ErrNoSession                          = NewError(E_NO_SESSION, E_MSG_NO_SESSION)
	ErrSprotoWaitPushTimeout              = NewError(E_SPROTO_WAIT_PUSH_TIMEOUT, E_MSG_SPROTO_WAIT_PUSH_TIMEOUT)
	ErrSprotoWaitPushConflict             = NewError(E_SPROTO_WAIT_PUSH_CONFLICT, E_MSG_SPROTO_WAIT_PUSH_CONFLICT)
	ErrTodo                               = NewError(E_TODO, E_MSG_TODO)
	ErrServerDown                         = NewError(E_SERVER_DOWN, E_MSG_SERVER_DOWN)
	ErrGrpcServerNotOpen                  = NewError(E_GRPC_SERVER_NOT_OPEN, E_MSG_GRPC_SERVER_NOT_OPEN)
	ErrAdminClientNotImplement            = NewError(E_ADMIN_CLIENT_NOT_IMPLEMENT, E_MSG_ADMIN_CLIENT_NOT_IMPLEMENT)
	ErrRegistryNOtImplement               = NewError(E_REGISTRY_NOT_IMPLEMENT, E_MSG_REGISTRY_NOT_IMPLEMENT)
	ErrInvalidMemberId                    = NewError(E_INVALID_MEMBER_ID, E_MSG_INVALID_MEMBER_ID)
	ErrBehaviorNotInit                    = NewError(E_BEHAVIOR_DRIVER_NOT_INIT, E_MSG_BEHAVIOR_DRIVER_NOT_INIT)
	ErrDbNotSupport                       = NewError(E_DB_NOT_SUPPORT, E_MSG_DB_NOT_SUPPORT)
	ErrDbNotConfig                        = NewError(E_DB_NOT_CONFIG, E_MSG_DB_NOT_CONFIG)
	ErrInvalidService                     = NewError(E_INVALID_SERVICE, E_MSG_INVALID_SERVICE)
	ErrServiceNotFound                    = NewError(E_SERVICE_NOT_FOUND, E_MSG_SERVICE_NOT_FOUND)
	ErrServiceLocked                      = NewError(E_SERVICE_LOCKED, E_MSG_SERVICE_LOCKED)
	ErrAlreadyDestory                     = NewError(E_ALREADY_DESTORY, E_MSG_ALREADY_DESTORY)
	ErrServiceContainerNotImplement       = NewError(E_SERVICE_CONTAINER_NOT_IMPLEMENT, E_MSG_SERVICE_ONTAINER_NOT_IMPLEMENT)
	ErrServiceNotImplement                = NewError(E_SERVICE_NOT_IMPLEMENT, E_MSG_SERVICE_NOT_IMPLEMENT)
	ErrUserNotFound                       = NewError(E_USER_NOT_FOUND, E_MSG_USER_NOT_FOUND)
	ErrActorCallTimeOut                   = NewError(E_ACTOR_CALL_TIME_OUT, E_MSG_ACTOR_CALL_TIME_OUT)
	ErrEtcdConfigNotFound                 = NewError(E_ETCD_CONFIG_NOT_FOUND, E_MSG_ETCD_CONFIG_NOT_FOUND)
	ErrGrpcConfigNotFound                 = NewError(E_GRPC_CONFIG_NOT_FOUND, E_MSG_GRPC_CONFIG_NOT_FOUND)
	ErrInterrupt                          = NewError(E_INTERRUPT, E_MSGS_INTERRUPT)
	ErrCatalogServerMetaNotFound          = NewError(E_CATALOG_SERVER_META_NOT_FOUND, E_MSG_CATALOG_SERVER_META_NOT_FOUND)
	ErrCatalogServerKeyNotFound           = NewError(E_CATALOG_SERVER_KEY_NOT_FOUND, E_MSG_CATALOG_SERVER_KEY_NOT_FOUND)
	ErrCatalogServerHandlerNotRegist      = NewError(E_CATALOG_SERVER_HANDLER_NOT_REGIST, E_MSG_CATALOG_SERVER_HANDLER_NOT_REGIST)
	ErrCatalogServerHandlerNotImplement   = NewError(E_CATALOG_SERVER_HANDLER_NOT_IMPLEMENT, E_MSG_CATALOG_SERVER_HANDLER_NOT_IMPLEMENT)
	ErrServiceAlreadyStopped              = NewError(E_SERVICE_ALREADY_STOPPED, E_MSG_SERVICE_ALREADY_STOPPED)
	ErrServiceAlreadyStarted              = NewError(E_SERVICE_ALREADY_STARTED, E_MSG_SERVICE_ALREADY_STARTED)
	ErrDbClientComponentNotImplement      = NewError(E_DB_CLIENT_COMPONENT_NOT_IMPLEMENT, E_MSG_DB_CLIENT_COMPONENT_NOT_IMPLEMENT)
	ErrServerNotFound                     = NewError(E_SERVER_NOT_FOUND, E_MSG_SERVER_NOT_FOUND)
	ErrSdkPayOrderCheckMethodNotImplement = NewError(E_SDK_PAY_ORDER_CHECK_METHOD_NOT_IMPLEMENT, E_MSG_SDK_PAY_ORDER_CHECK_METHOD_NOT_IMPLEMENT)
	ErrSdkPayOrderCheckArgsInvalid        = NewError(E_SDK_PAY_ORDER_CHECK_ARGS_INVALID, E_MSG_SDK_PAY_ORDER_CHECK_ARGS_INVALID)
	ErrAccountPlatformNotSupport          = NewError(E_ACCOUNT_PLATFORM_NOT_SUPPORT, E_MSG_ACCOUNT_PLATFORM_NOT_SUPPORT)
	ErrPayOrderStatusInvalid              = NewError(E_PAY_ORDER_STATUS_INVALID, E_MSG_PAY_ORDER_STATUS_INVALID)
	ErrPayOrderAmountInvalid              = NewError(E_PAY_ORDER_AMOUNT_INVALID, E_MSG_PAY_ORDER_AMOUNT_INVALID)
	ErrPayOrderAccountPlatformInvalid     = NewError(E_PAY_ORDER_ACCOUNT_PLATFORM_INVALID, E_MSG_PAY_ORDER_ACCOUNT_PLATFORM_INVALID)
	ErrCast                               = NewError(E_CAST, E_MSG_CAST)
	ErrProtoRequestCast                   = NewError(E_PROTO_REQUEST_CAST, E_MSG_PROTO_REQUEST_CAST)
	ErrProtoResponseCast                  = NewError(E_PROTO_RESPONSE_CAST, E_MSG_PROTO_RESPONSE_CAST)
	ErrProtoPushCast                      = NewError(E_PROTO_PUSH_CAST, E_MSG_PROTO_PUSH_CAST)
	ErrProtoResponseNew                   = NewError(E_PROTO_RESPONSE_NEW, E_MSG_PROTO_RESPONSE_NEW)
	ErrCronNotImplement                   = NewError(E_CRON_NOT_IMPLEMENT, E_MSG_CRON_NOT_IMPLEMENT)
	ErrConfigEnvInvalidSyntax             = NewError(E_CONFIG_ENV_INVALID_SYNTAX, E_MSG_CONFIG_ENV_INVALID_SYNTAX)
	ErrInvalidSyntax                      = NewError(E_INVALID_SYNTAX, E_MSG_INVALID_SYNTAX)
)

// 根据code, msg创建error
func NewError(code int32, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

// 根据code, msg, values创建error
func NewErrorw(code int32, msg string, values ...interface{}) *Error {
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
		Code:   code,
		Msg:    msg,
		Values: kvs,
	}
}

// 转换或者创建error
func NewOrCastError(err error) *Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	} else {
		msg := err.Error()
		arr := strings.Split(msg, ":")
		if len(arr) < 2 {
			return &Error{
				Code: E_MALFORMED,
				Msg:  msg,
			}
		} else if v, err := strconv.Atoi(arr[0]); err != nil {
			return &Error{
				Code: E_MALFORMED,
				Msg:  msg,
			}
		} else {
			return &Error{
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
	if e, ok := err.(*Error); ok {
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
	if e, ok := err.(*Error); ok {
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
	if e, ok := err.(*Error); ok {
		return string(e.Stack)
	}
	return ""
}

type Error struct {
	Code   int32
	Msg    string
	Stack  []byte
	Values map[string]interface{}
}

// 拷贝error并且附加stack
func (e *Error) Trace() *Error {
	return &Error{Code: e.Code, Msg: e.Msg, Stack: debug.Stack(), Values: e.Values}
}

// 格式化错误字符串
func (e *Error) Error() string {
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
func (err *Error) WithValues(values ...interface{}) *Error {
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
func (err *Error) WithTrace(values ...interface{}) *Error {
	err.Stack = debug.Stack()
	return err
}

// 附加错误描述到error中
func (err *Error) WithMsg(msg string) *Error {
	err.Msg = msg
	return err
}

func (err *Error) WithFile(file string) *Error {
	if err.Values == nil {
		err.Values = make(map[string]interface{})
	}
	err.Values["file"] = file
	return err
}

func (err *Error) WithLines(content []byte) *Error {
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
