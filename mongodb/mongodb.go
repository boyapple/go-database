package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(name string, opts ...Option) (*mongo.Client, error) {
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	ctx := context.Background()
	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(opt.Uri))
	if err != nil {
		return nil, err
	}
	if err = cli.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return cli, nil
}
