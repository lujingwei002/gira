package service_options

// ====== register options ===================

// app内唯一的服务
func WithAsAppServiceOption(enabled bool) AsAppServiceOption {
	return AsAppServiceOption{
		enabled: enabled,
	}
}

// ====== where options ===================
// 最大数量
func WithWhereMaxCountOption(count int) WhereMaxCountOption {
	return WhereMaxCountOption{
		count: count,
	}
}

// 设置正则表达式查找
func WithWhereRegexOption(regex string) WhereRegexOption {
	return WhereRegexOption{
		regex: regex,
	}
}

// ====== where options ===================
type WhereOptions struct {
	MaxCount int
	Regex    string
}

type WhereOption interface {
	ConfigWhereOption(opts *WhereOptions)
}

type WhereMaxCountOption struct {
	count int
}

func (opt WhereMaxCountOption) ConfigWhereOption(opts *WhereOptions) {
	opts.MaxCount = opt.count
}

type WhereRegexOption struct {
	regex string
}

func (opt WhereRegexOption) ConfigWhereOption(opts *WhereOptions) {
	opts.Regex = opt.regex
}

// ====== register options ===================

type RegisterOptions struct {
	AsAppService bool
}

type RegisterOption interface {
	ConfigRegisterOption(opts *RegisterOptions)
}

type AsAppServiceOption struct {
	enabled bool
}

func (opt AsAppServiceOption) ConfigRegisterOption(opts *RegisterOptions) {
	opts.AsAppService = opt.enabled
}
