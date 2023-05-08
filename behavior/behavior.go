package behavior

import (
	"sync/atomic"
	"time"

	"github.com/lujingwei002/gira/facade"
)

type SyncOptions struct {
	Step int
}

type SyncOption interface {
	ConfigSyncOptions(opts *SyncOptions)
}

type BatchInsertOption struct {
	step int
}

func (opt BatchInsertOption) ConfigSyncOptions(opts *SyncOptions) {
	opts.Step = opt.step
}

func WitchBatchInsertOption(step int) BatchInsertOption {
	return BatchInsertOption{
		step: step,
	}
}

var seqId int64

func GetSeqId() int64 {
	seqId = atomic.AddInt64(&seqId, 1)
	return time.Now().Unix()<<32 + int64(facade.GetAppId())<<24 + (seqId & 0xffffff)
}

type MigrateOptions struct {
	EnabledDropIndex bool
}
type MigrateOption interface {
	ConfigMigrateOptions(opts *MigrateOptions)
}

type MigrateDropIndex struct {
	enabled bool
}

func (opt MigrateDropIndex) ConfigMigrateOptions(opts *MigrateOptions) {
	opts.EnabledDropIndex = opt.enabled
}

func WithMigrateDropIndex(enabled bool) MigrateDropIndex {
	return MigrateDropIndex{
		enabled: enabled,
	}
}
