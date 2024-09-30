package db

import "github.com/boyapple/go-database/gorm/condition"

type Options struct {
	MultiCondition condition.MultiCondition
}

type Option func(*Options)

func WithCondition(conditions ...condition.Condition) Option {
	return func(o *Options) {
		o.MultiCondition = conditions
	}
}
