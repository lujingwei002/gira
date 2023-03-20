package macro

type SyncRequest interface {
	SyncCall()
}
type SyncBase struct {
	__sync_ch__ chan SyncRequest
	SyncLock__  int64
}

func NewSyncBase(size int) *SyncBase {
	return &SyncBase{
		__sync_ch__: make(chan SyncRequest, size),
	}
}
func (self *SyncBase) SyncChan() chan SyncRequest {
	return self.__sync_ch__
}
