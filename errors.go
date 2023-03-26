package gira

import (
	"strconv"
	"strings"
)

// 错误码常量
const (
	E_OK                             = 0
	E_MALFORMED                      = -1 // error格式错误,不是"code:msg"格式
	E_UNKNOWN                        = -2
	E_NULL_POINTER                   = -2
	E_RESOURCE_MANAGER_NOT_IMPLEMENT = -3
	E_CONFIG_HANDLER_NOT_IMPLEMENT   = -4
	E_RESOURCE_LOADER_NOT_IMPLEMENT  = -5
	E_RESOURCE_HANDLER_NOT_IMPLEMENT = -6
	E_NULL_OBJECT                    = -7
	E_HTTP_HANDLER_NOT_IMPLEMENT     = -8
	E_REGISTER_SERVER_FAIL           = -9
	E_INVALID_PEER                   = -10
	E_DATA_UPSERT_FAIL               = -11
	E_DATA_EXIST                     = -12
	E_DATA_NOT_EXIST                 = -13
	E_DATA_NOT_FOUND                 = -14
	E_DATA_INSERT_FAIL               = -15
	E_DATA_DELETE_FAIL               = -16
	E_SDK_NOT_IMPLEMENT              = -17
	E_INVALID_SDK_TOKEN              = -18
	E_INVALID_ARGS                   = -19
	E_INVALID_JWT                    = -20
	E_JWT_EXPIRE                     = -21
	E_GATE_HANDLER_NOT_IMPLEMENT     = -22
	E_PEER_HANDLER_NOT_IMPLEMENT     = -23
	E_GRPC_HANDLER_NOT_IMPLEMENT     = -24
	E_HALL_HANDLER_NOT_IMPLEMENT     = -25
	E_SPROTO_HANDLER_NOT_IMPLEMENT   = -26
	E_SPROTO_REQ_ID_CONFLICT         = -27
	E_SPROTO_REQ_TIMEOUT             = -28
	E_SPROTO_RESPONSE_CONVERSION     = -29
	E_READ_ON_CLOSED_CLIENT          = -30
	E_PEER_NOT_FOUND                 = -31
	E_MEMBER_INSTEAD                 = -32
	E_MEMBER_LOCKED                  = -33
	E_GRPC_CLIENT_POOL_NIL           = -34
	E_BROKEN_CHANNEL                 = -35
	E_SESSION_CLOSED                 = -36
	E_UPSTREAM_UNAVAILABLE           = -37
	E_UPSTREAM_UNREACHABLE           = -38
	E_SERVICE_UNAVAILABLE            = -39
	E_SPROTO_PUSH_CONVERSION         = -40
	E_PROJECT_FILE_NOT_FOUND         = -41
	E_INVALID_PASSWORD               = -42
	E_GEN_NOT_CHANGE                 = -43
)
const (
	E_MSG_OK                             = "成功"
	E_MSG_MALFORMED                      = "malformed error"
	E_MSG_UNKNOWN                        = "未知错误"
	E_MSG_NULL_POINTER                   = "空指针"
	E_MSG_RESOURCE_MANAGER_NOT_IMPLEMENT = "resource manager not implement"
	E_MSG_CONFIG_HANDLER_NOT_IMPLEMENT   = "config handler not implement"
	E_MSG_RESOURCE_LOADER_NOT_IMPLEMENT  = "resource loader not implement"
	E_MSG_RESOURCE_HANDLER_NOT_IMPLEMENT = "resource handler not implement"
	E_MSG_NULL_OBJECT                    = "null object"
	E_MSG_HTTP_HANDLER_NOT_IMPLEMENT     = "http handler not implement"
	E_MSG_REGISTER_SERVER_FAIL           = "register server fail"
	E_MSG_INVALID_PEER                   = "invalid peer"
	E_MSG_DATA_UPSERT_FAIL               = "upsert fail"
	E_MSG_DATA_EXIST                     = "data exist"
	E_MSG_DATA_NOT_EXIST                 = "data not exist"
	E_MSG_DATA_NOT_FOUND                 = "data not found"
	E_MSG_DATA_INSERT_FAIL               = "data insert fail"
	E_MSG_DATA_DELETE_FAIL               = "data delete fail"
	E_MSG_SDK_NOT_IMPLEMENT              = "无法识别的sdk"
	E_MSG_INVALID_SDK_TOKEN              = "无效的sdk token"
	E_MSG_INVALID_ARGS                   = "无效参数"
	E_MSG_INVALID_JWT                    = "无效的token"
	E_MSG_JWT_EXPIRE                     = "token已过期"
	E_MSG_GATE_HANDLER_NOT_IMPLEMENT     = "gate handler not implement"
	E_MSG_PEER_HANDLER_NOT_IMPLEMENT     = "peer handler not implement"
	E_MSG_GRPC_HANDLER_NOT_IMPLEMENT     = "grpc handler not implement"
	E_MSG_HALL_HANDLER_NOT_IMPLEMENT     = "hall handler not implement"
	E_MSG_SPROTO_HANDLER_NOT_IMPLEMENT   = "sproto handler not implement"
	E_MSG_SPROTO_REQ_ID_CONFLICT         = "sproto req id conflict"
	E_MSG_SPROTO_REQ_TIMEOUT             = "sproto req id timeout"
	E_MSG_SPROTO_RESPONSE_CONVERSION     = "sproto response type conversion"
	E_MSG_READ_ON_CLOSED_CLIENT          = "read on closed client"
	E_MSG_PEER_NOT_FOUND                 = "peer not found"
	E_MSG_MEMBER_INSTEAD                 = "账号在其他地方登录"
	E_MSG_MEMBER_LOCKED                  = "账号在其他地方被锁定"
	E_MSG_GRPC_CLIENT_POOL_NIL           = "grpc pool无法申请client"
	E_MSG_BROKEN_CHANNEL                 = "管道已关闭，不能再写数据"
	E_MSG_SESSION_CLOSED                 = "会话已经关闭"
	E_MSG_UPSTREAM_UNAVAILABLE           = "上游服务不可用"
	E_MSG_UPSTREAM_UNREACHABLE           = "上游服务不可达"
	E_MSG_SERVICE_UNAVAILABLE            = "服务不可用"
	E_MSG_SPROTO_PUSH_CONVERSION         = "sproto push type conversion"
	E_MSG_PROJECT_FILE_NOT_FOUND         = "gira.yaml文件找不到"
	E_MSG_INVALID_PASSWORD               = "密码错误"
	E_MSG_GEN_NOT_CHANGE                 = "gen源文件没变化"
)

type Error struct {
	Code int32
	Msg  string
}

func (e *Error) Error() string {
	return e.Msg
}
func NewError(code int32, msg string) *Error {
	return &Error{
		code,
		msg,
	}
}

var (
	ErrNullPonter                  = NewError(E_RESOURCE_MANAGER_NOT_IMPLEMENT, E_MSG_NULL_POINTER)
	ErrMalformed                   = NewError(E_MALFORMED, E_MSG_MALFORMED)
	ErrUnknown                     = NewError(E_UNKNOWN, E_MSG_UNKNOWN)
	ErrResourceManagerNotImplement = NewError(E_CONFIG_HANDLER_NOT_IMPLEMENT, E_MSG_RESOURCE_MANAGER_NOT_IMPLEMENT)
	ErrConfigHandlerNotImplement   = NewError(E_CONFIG_HANDLER_NOT_IMPLEMENT, E_MSG_CONFIG_HANDLER_NOT_IMPLEMENT)
	ErrResourceLoaderNotImplement  = NewError(E_RESOURCE_LOADER_NOT_IMPLEMENT, E_MSG_RESOURCE_LOADER_NOT_IMPLEMENT)
	ErrResourceHandlerNotImplement = NewError(E_RESOURCE_HANDLER_NOT_IMPLEMENT, E_MSG_RESOURCE_HANDLER_NOT_IMPLEMENT)
	ErrNullObject                  = NewError(E_NULL_OBJECT, E_MSG_NULL_OBJECT)
	ErrHttpHandlerNotImplement     = NewError(E_HTTP_HANDLER_NOT_IMPLEMENT, E_MSG_HTTP_HANDLER_NOT_IMPLEMENT)
	ErrRegisterServerFail          = NewError(E_REGISTER_SERVER_FAIL, E_MSG_REGISTER_SERVER_FAIL)
	ErrInvalidPeer                 = NewError(E_INVALID_PEER, E_MSG_INVALID_PEER)
	ErrDataUpsertFail              = NewError(E_DATA_UPSERT_FAIL, E_MSG_DATA_UPSERT_FAIL)
	ErrDataExist                   = NewError(E_DATA_EXIST, E_MSG_DATA_EXIST)
	ErrDataNotExist                = NewError(E_DATA_NOT_EXIST, E_MSG_DATA_NOT_EXIST)
	ErrDataNotFound                = NewError(E_DATA_NOT_FOUND, E_MSG_DATA_NOT_FOUND)
	ErrDataInsertFail              = NewError(E_DATA_INSERT_FAIL, E_MSG_DATA_INSERT_FAIL)
	ErrDataDeleteFail              = NewError(E_DATA_DELETE_FAIL, E_MSG_DATA_DELETE_FAIL)
	ErrSdkNotImplement             = NewError(E_SDK_NOT_IMPLEMENT, E_MSG_SDK_NOT_IMPLEMENT)
	ErrInvalidSdkToken             = NewError(E_INVALID_SDK_TOKEN, E_MSG_INVALID_SDK_TOKEN)
	ErrInvalidArgs                 = NewError(E_INVALID_ARGS, E_MSG_INVALID_ARGS)
	ErrInvalidJwt                  = NewError(E_INVALID_JWT, E_MSG_INVALID_JWT)
	ErrJwtExpire                   = NewError(E_JWT_EXPIRE, E_MSG_JWT_EXPIRE)
	ErrGateHandlerNotImplement     = NewError(E_GATE_HANDLER_NOT_IMPLEMENT, E_MSG_GATE_HANDLER_NOT_IMPLEMENT)
	ErrPeerHandlerNotImplement     = NewError(E_PEER_HANDLER_NOT_IMPLEMENT, E_MSG_PEER_HANDLER_NOT_IMPLEMENT)
	ErrGrpcHandlerNotImplement     = NewError(E_GRPC_HANDLER_NOT_IMPLEMENT, E_MSG_GRPC_HANDLER_NOT_IMPLEMENT)
	ErrHallHandlerNotImplement     = NewError(E_HALL_HANDLER_NOT_IMPLEMENT, E_MSG_HALL_HANDLER_NOT_IMPLEMENT)
	ErrSprotoHandlerNotImplement   = NewError(E_SPROTO_HANDLER_NOT_IMPLEMENT, E_MSG_SPROTO_HANDLER_NOT_IMPLEMENT)
	ErrSprotoReqIdConflict         = NewError(E_SPROTO_REQ_ID_CONFLICT, E_MSG_SPROTO_REQ_ID_CONFLICT)
	ErrSprotoReqTimeout            = NewError(E_SPROTO_REQ_TIMEOUT, E_MSG_SPROTO_REQ_TIMEOUT)
	ErrSprotoResponseConversion    = NewError(E_SPROTO_RESPONSE_CONVERSION, E_MSG_SPROTO_RESPONSE_CONVERSION)
	ErrReadOnClosedClient          = NewError(E_READ_ON_CLOSED_CLIENT, E_MSG_READ_ON_CLOSED_CLIENT)
	ErrPeerNotFound                = NewError(E_PEER_NOT_FOUND, E_MSG_PEER_NOT_FOUND)
	ErrMemberInstead               = NewError(E_MEMBER_INSTEAD, E_MSG_MEMBER_INSTEAD)
	ErrMemberLocked                = NewError(E_MEMBER_LOCKED, E_MSG_MEMBER_LOCKED)
	ErrGrpcClientPoolNil           = NewError(E_GRPC_CLIENT_POOL_NIL, E_MSG_GRPC_CLIENT_POOL_NIL)
	ErrBrokenChannel               = NewError(E_BROKEN_CHANNEL, E_MSG_BROKEN_CHANNEL)
	ErrSessionClosed               = NewError(E_SESSION_CLOSED, E_MSG_SESSION_CLOSED)
	ErrUpstreamUnavailable         = NewError(E_UPSTREAM_UNAVAILABLE, E_MSG_UPSTREAM_UNAVAILABLE)
	ErrUpstreamUnreachable         = NewError(E_UPSTREAM_UNREACHABLE, E_MSG_UPSTREAM_UNREACHABLE)
	ErrServiceUnavailable          = NewError(E_SERVICE_UNAVAILABLE, E_MSG_SERVICE_UNAVAILABLE)
	ErrSprotoPushConversion        = NewError(E_SPROTO_PUSH_CONVERSION, E_MSG_SPROTO_PUSH_CONVERSION)
	ErrProjectFileNotFound         = NewError(E_PROJECT_FILE_NOT_FOUND, E_MSG_PROJECT_FILE_NOT_FOUND)
	ErrInvalidPassword             = NewError(E_INVALID_PASSWORD, E_MSG_INVALID_PASSWORD)
	ErrGenNotChange                = NewError(E_GEN_NOT_CHANGE, E_MSG_GEN_NOT_CHANGE)
)

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
