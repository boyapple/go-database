package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func New(name string, opts ...Option) (redis.UniversalClient, error) {
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	cli := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    opt.Addrs,
		Password: opt.Password,
		DB:       opt.DB,
	})
	if err := cli.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return cli, nil
}
