package redis

type Options struct {
	Addrs    []string
	Password string
	DB       int
}

type Option func(*Options)

func WithAddrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

func WithPassword(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

func WithDB(db int) Option {
	return func(o *Options) {
		o.DB = db
	}
}
