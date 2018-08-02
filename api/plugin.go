package api

type Plugin interface {
	Init(options *Options) error
}
