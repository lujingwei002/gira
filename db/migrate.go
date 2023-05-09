package db

type MigrateOptions struct {
	EnabledDropIndex bool
	ConnectTimeout   int64
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

type MigrateConnectTimeout struct {
	timeout int64
}

func (opt MigrateConnectTimeout) ConfigMigrateOptions(opts *MigrateOptions) {
	opts.ConnectTimeout = opt.timeout
}

func WithMigrateConnectTimeout(timeout int64) MigrateConnectTimeout {
	return MigrateConnectTimeout{
		timeout: timeout,
	}
}
