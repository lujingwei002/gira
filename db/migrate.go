package db

// ========= migrate选项 ==================
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

// 设置允许drop index
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

// 设置连接超时时间
func WithMigrateConnectTimeout(timeout int64) MigrateConnectTimeout {
	return MigrateConnectTimeout{
		timeout: timeout,
	}
}
