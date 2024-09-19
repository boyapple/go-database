package es

import (
	"context"

	"github.com/olivere/elastic/v7"
)

func New(opts ...Option) (*elastic.Client, error) {
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	cli, err := elastic.NewClient(
		elastic.SetURL(opt.Url),
		elastic.SetSniff(true),
	)
	if err != nil {
		return nil, err
	}
	_, _, err = cli.Ping(opt.Url).Do(context.Background())
	if err != nil {
		return nil, err
	}
	return cli, nil
}
