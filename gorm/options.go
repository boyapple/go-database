package gorm

type Options struct {
	Dsn string
}

type Option func(*Options)

func WithDsn(dsn string) Option {
	return func(o *Options) {
		o.Dsn = dsn
	}
}
