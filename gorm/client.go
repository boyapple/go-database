package db

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Client 只用于单个对象(单张表)的client,传入的类型T必须实现 TableName()string 函数
type Client[T any] interface {
	Get(ctx context.Context, condition Condition) (T, error)
	List(ctx context.Context, opts ...Option) ([]T, error)
	Count(ctx context.Context, opts ...Option) (int64, error)
	Create(ctx context.Context, t T, opts ...Option) error
	Updates(ctx context.Context, t T, opts ...Option) error
}

func NewClient[T any](serviceName string) Client[T] {
	return &gormClient[T]{serviceName: serviceName}
}

type gormClient[T any] struct {
	serviceName string
}

func (c *gormClient[T]) Get(ctx context.Context, condition Condition) (T, error) {
	var t T
	if condition == nil {
		return t, fmt.Errorf("must setup condition")
	}
	db, err := c.getDB(ctx)
	if err != nil {
		return t, err
	}
	where, err := condition.Where()
	if err != nil {
		return t, err
	}
	if err = db.Scopes(where).First(&t).Error; err != nil {
		return t, err
	}
	return t, nil
}

func (c *gormClient[T]) List(ctx context.Context, opts ...Option) ([]T, error) {
	db, err := c.getDB(ctx)
	if err != nil {
		return nil, err
	}
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	if len(opt.Condition) > 0 {
		scopes, err := opt.Condition.Build()
		if err != nil {
			return nil, err
		}
		db.Scopes(scopes...)
	}
	if opt.Page != nil {
		if err = db.Count(&opt.Page.Count).Error; err != nil {
			return nil, err
		}
		db.Scopes(func(db *gorm.DB) *gorm.DB {
			offset := opt.Page.GetOffset()
			limit := opt.Page.GetLimit()
			return db.Offset((offset - 1) * limit).Limit(limit)
		})
	}
	var list []T
	if err = db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (c *gormClient[T]) Count(ctx context.Context, opts ...Option) (int64, error) {
	db, err := c.getDB(ctx)
	if err != nil {
		return 0, err
	}
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	if len(opt.Condition) > 0 {
		scopes, err := opt.Condition.Build()
		if err != nil {
			return 0, err
		}
		db.Scopes(scopes...)
	}
	var count int64
	if err = db.Count(&count).Error; err != nil {
		return 0, nil
	}
	return count, nil
}

func (c *gormClient[T]) Create(ctx context.Context, t T, opts ...Option) error {
	db, err := c.getDB(ctx)
	if err != nil {
		return err
	}
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	if len(opt.OnlyColumn) > 0 {
		var columns []clause.Column
		for _, column := range opt.OnlyColumn {
			columns = append(columns, clause.Column{Name: column})
		}
		conflict := &clause.OnConflict{
			Columns: columns,
		}
		if len(opt.UpdateColumn) == 0 {
			conflict.UpdateAll = true
		} else {
			conflict.DoUpdates = clause.AssignmentColumns(opt.UpdateColumn)
		}
		db.Clauses(conflict)
	}
	return db.Create(t).Error
}

func (c *gormClient[T]) Updates(ctx context.Context, t T, opts ...Option) error {
	db, err := c.getDB(ctx)
	if err != nil {
		return err
	}
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	if len(opt.UpdateColumn) > 0 {
		db.Select(opt.UpdateColumn)
	}
	return db.Updates(t).Error
}

func (c *gormClient[T]) getDB(ctx context.Context) (*gorm.DB, error) {
	db, err := Get(c.serviceName)
	if err != nil {
		return nil, err
	}
	var t T
	return db.WithContext(ctx).Model(&t), nil
}
