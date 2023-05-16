package options

// where option
type WhereOptions struct {
	MulticastCount int
}

type WhereOption interface {
	ConfigWhereOption(opts *WhereOptions)
}

type WhereMulticastOption struct {
	count int
}

func WithWhereMulticastOption(count int) WhereMulticastOption {
	return WhereMulticastOption{
		count: count,
	}
}

func (opt WhereMulticastOption) ConfigWhereOption(opts *WhereOptions) {
	opts.MulticastCount = opt.count
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

func WithRegisterAsGroupOption(asGroup bool) RegisterAsGroupOption {
	return RegisterAsGroupOption{
		asGroup: asGroup,
	}
}

func (opt RegisterAsGroupOption) ConfigRegisterOption(opts *RegisterOptions) {
	opts.AsGroup = opt.asGroup
}

type RegisterCatAppidOption struct {
	catAppId bool
}

func WithRegisterCatAppIdOption(catAppId bool) RegisterCatAppidOption {
	return RegisterCatAppidOption{
		catAppId: catAppId,
	}
}

func (opt RegisterCatAppidOption) ConfigRegisterOption(opts *RegisterOptions) {
	opts.CatAppId = opt.catAppId
}
