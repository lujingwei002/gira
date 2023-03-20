package actor

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
