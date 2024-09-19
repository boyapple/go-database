package es

type Options struct {
	Url      string
	Username string
	Password string
}

type Option func(*Options)

func WithUrl(url string) Option {
	return func(o *Options) {
		o.Url = url
	}
}

func WithUsername(username string) Option {
	return func(o *Options) {
		o.Username = username
	}
}

func WithPassword(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}
