package registry_options

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

// 注册成组员
func WithRegisterAsGroupOption(asGroup bool) RegisterAsGroupOption {
	return RegisterAsGroupOption{
		asGroup: asGroup,
	}
}

// 当注册成组员时，使用app id作为组员id
func WithRegisterCatAppIdOption(catAppId bool) RegisterCatAppidOption {
	return RegisterCatAppidOption{
		catAppId: catAppId,
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

type RegisterOptions struct {
	AsGroup  bool
	CatAppId bool
}

type RegisterOption interface {
	ConfigRegisterOption(opts *RegisterOptions)
}

type RegisterAsGroupOption struct {
	asGroup bool
}

func (opt RegisterAsGroupOption) ConfigRegisterOption(opts *RegisterOptions) {
	opts.AsGroup = opt.asGroup
}

type RegisterCatAppidOption struct {
	catAppId bool
}

func (opt RegisterCatAppidOption) ConfigRegisterOption(opts *RegisterOptions) {
	opts.CatAppId = opt.catAppId
}
