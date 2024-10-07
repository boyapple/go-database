package db

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Client 只支持单表操作的gorm client
type Client[T schema.Tabler] interface {
	// Get 获取单条数据,至少设置一个条件
	Get(ctx context.Context, opts ...Option) (T, error)
	// List 获取列表数据
	List(ctx context.Context, opts ...Option) ([]T, error)
	// Count 获取总数
	Count(ctx context.Context, opts ...Option) (int64, error)
	// Create 创建数据,支持on duplicate key update
	Create(ctx context.Context, t T, opts ...Option) error
	// Update 根据结构体更新,没有条件会根据表的主键id更新
	Update(ctx context.Context, t T, opts ...Option) error
	// Updates 根据键值对更新,至少设置一个条件
	Updates(ctx context.Context, keyValue map[string]interface{}, opts ...Option) error
}

func New[T schema.Tabler](serviceName string) Client[T] {
	var t T
	return &impl[T]{
		serviceName: serviceName,
		tableName:   t.TableName(),
	}
}

type impl[T any] struct {
	serviceName string
	tableName   string
}

func (i *impl[T]) Get(ctx context.Context, opts ...Option) (T, error) {
	var t T
	db, err := i.getDB(ctx)
	if err != nil {
		return t, err
	}
	opt := i.getOptions(opts...)
	if len(opt.Columns) > 0 {
		db.Select(opt.Columns)
	}
	if len(opt.MultiCondition) == 0 {
		return t, fmt.Errorf("least contion one condition")
	}
	conditions, err := opt.MultiCondition.Build()
	if err != nil {
		return t, err
	}
	if err = db.Scopes(conditions...).First(&t).Error; err != nil {
		return t, err
	}
	return t, nil
}

func (i *impl[T]) List(ctx context.Context, opts ...Option) ([]T, error) {
	db, err := i.getDB(ctx)
	if err != nil {
		return nil, err
	}
	opt := i.getOptions(opts...)
	if len(opt.Columns) > 0 {
		db.Select(opt.Columns)
	}
	if len(opt.MultiCondition) > 0 {
		conditions, err := opt.MultiCondition.Build()
		if err != nil {
			return nil, err
		}
		db.Scopes(conditions...)
	}
	var list []T
	if err = db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (i *impl[T]) Count(ctx context.Context, opts ...Option) (int64, error) {
	db, err := i.getDB(ctx)
	if err != nil {
		return 0, err
	}
	opt := i.getOptions(opts...)
	if len(opt.MultiCondition) > 0 {
		conditions, err := opt.MultiCondition.Build()
		if err != nil {
			return 0, err
		}
		db.Scopes(conditions...)
	}
	var count int64
	if err = db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (i *impl[T]) Create(ctx context.Context, t T, opts ...Option) error {
	opt := i.getOptions(opts...)
	db, err := i.getDB(ctx)
	if err != nil {
		return err
	}
	if len(opt.MultiCondition) > 0 {
		cond, err := opt.MultiCondition[0].Compile()
		if err != nil {
			return err
		}
		db.Scopes(cond)
	}
	return db.Create(t).Error
}

func (i *impl[T]) Update(ctx context.Context, t T, opts ...Option) error {
	db, err := i.getDB(ctx)
	if err != nil {
		return err
	}
	opt := i.getOptions(opts...)
	if len(opt.Columns) > 0 {
		db.Select(opt.Columns)
	}
	if len(opt.MultiCondition) > 0 {
		conditions, err := opt.MultiCondition.Build()
		if err != nil {
			return err
		}
		db.Scopes(conditions...)
	}
	return db.Updates(t).Error
}

func (i *impl[T]) Updates(ctx context.Context, keyValue map[string]interface{}, opts ...Option) error {
	db, err := i.getDB(ctx)
	if err != nil {
		return err
	}
	opt := i.getOptions(opts...)
	if len(opt.MultiCondition) == 0 {
		return fmt.Errorf("least contion one condition")
	}
	conditions, err := opt.MultiCondition.Build()
	if err != nil {
		return err
	}
	return db.Scopes(conditions...).Updates(keyValue).Error
}

func (i *impl[T]) getDB(ctx context.Context) (*gorm.DB, error) {
	db, err := dbMux.Get(i.serviceName)
	if err != nil {
		return nil, err
	}
	return db.WithContext(ctx).Table(i.tableName), nil
}

func (i *impl[T]) getOptions(opts ...Option) *Options {
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	return opt
}
