package actor

import (
	"time"

	"github.com/lujingwei002/gira"
)

var ErrCallTimeOut = gira.NewError(-1, "call time out")

type Request interface {
	Call()
}

type Actor struct {
	__sync_ch__ chan Request
	SyncLock__  int64
}

func NewActor(size int) *Actor {
	return &Actor{
		__sync_ch__: make(chan Request, size),
	}
}

func (self *Actor) Inbox() chan Request {
	return self.__sync_ch__
}

type CallOptions struct {
	TimeOut time.Duration
}
type CallTimeOutOption time.Duration

func WithCallTimeOut(timeout time.Duration) CallTimeOutOption {
	return CallTimeOutOption(timeout)
}
func (self CallTimeOutOption) Config(opt *CallOptions) {
	opt.TimeOut = time.Duration(self)
}

type CallOption interface {
	Config(opt *CallOptions)
}
