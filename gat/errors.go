package gat

import "errors"

// Errors that could be occurred during message handling.
var (
	ErrSessionOnNotify    = errors.New("current session working on notify mode")
	ErrCloseClosedSession = errors.New("close closed session")
	ErrInvalidRegisterReq = errors.New("invalid register request")
	// ErrBrokenPipe represents the low-level connection has broken.
	ErrBrokenPipe = errors.New("broken low-level pipe")
	// ErrBufferExceed indicates that the current session buffer is full and
	// can not receive more data.
	ErrBufferExceed       = errors.New("session send buffer exceed")
	ErrCloseClosedGroup   = errors.New("close closed group")
	ErrClosedGroup        = errors.New("group closed")
	ErrMemberNotFound     = errors.New("member not found in the group")
	ErrSessionDuplication = errors.New("session has existed in the current group")
	ErrSprotoRequestType  = errors.New("sproto request type")
	ErrSprotoResponseType = errors.New("sproto response type")
	ErrHandShake          = errors.New("handshake failed")
	ErrHeartbeatTimeout   = errors.New("heartbeat timeout")
	ErrDialTimeout        = errors.New("dial timeout")
	ErrDialInterrupt      = errors.New("dial interrupt")
	ErrInvalidAddress     = errors.New("invalid address")
	ErrInvalidMessage     = errors.New("invalid message")
	ErrInvalidHandler     = errors.New("invalid handler")
	ErrInvalidPacket      = errors.New("invalid packet")
	ErrNotHandshake       = errors.New("没握手成功就发了数据包过来")
	ErrNotWorking         = errors.New("非正常状态，无法发送消息")
	ErrInvalidState       = errors.New("invalid state")
	ErrKick               = errors.New("被踢下线")
	ErrConnClosed         = errors.New("连接已关闭")
	ErrConnNotReady       = errors.New("连接末准备好")
	ErrHandShakeAck       = errors.New("handshake ack 出错")
)
