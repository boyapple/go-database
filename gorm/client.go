package gorm

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Client[T any] interface {
	Get(ctx context.Context, condition Condition) (T, error)
	List(ctx context.Context, opts ...Option) ([]T, error)
	Count(ctx context.Context, opts ...Option) (int64, error)
	Create(ctx context.Context, t T) error
	Update(ctx context.Context, t T) error
}

func NewClient[T any](name string, cfgOpts ...ConfigOption) (Client[T], error) {
	cfg := &Config{}
	for _, o := range cfgOpts {
		o(cfg)
	}
	db, err := New(name, cfg)
	if err != nil {
		return nil, err
	}
	return &client[T]{db: db}, nil
}

type client[T any] struct {
	db *gorm.DB
}

func (c *client[T]) Get(ctx context.Context, condition Condition) (T, error) {
	var t T
	if condition == nil {
		return t, fmt.Errorf("must setup condition")
	}
	if err := c.db.WithContext(ctx).Scopes(condition.Where()).First(&t).Error; err != nil {
		return t, err
	}
	return t, nil
}

func (c *client[T]) List(ctx context.Context, opts ...Option) ([]T, error) {
	var list []T
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	db := c.db.WithContext(ctx)
	if opt.Condition != nil {
		db.Scopes(opt.Condition.Where())
	}
	if opt.Page != nil {
		db.Scopes(func(db *gorm.DB) *gorm.DB {
			offset := opt.Page.GetOffset()
			limit := opt.Page.GetLimit()
			return db.Offset((offset - 1) * limit).Limit(limit)
		})
	}
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (c *client[T]) Count(ctx context.Context, opts ...Option) (int64, error) {
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	var t T
	db := c.db.WithContext(ctx).Model(&t)
	if opt.Condition != nil {
		db.Scopes(opt.Condition.Where())
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, nil
	}
	return count, nil
}

func (c *client[T]) Create(ctx context.Context, t T) error {
	return c.db.WithContext(ctx).Create(t).Error
}

func (c *client[T]) Update(ctx context.Context, t T) error {
	return c.db.WithContext(ctx).Updates(t).Error
}
