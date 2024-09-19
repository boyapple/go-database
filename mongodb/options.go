package mongodb

type Options struct {
	Uri string
}

type Option func(*Options)

func WithUri(uri string) Option {
	return func(o *Options) {
		o.Uri = uri
	}
}
