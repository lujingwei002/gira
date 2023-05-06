package behavior

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
