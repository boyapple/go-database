package db

import (
	"context"
	"github.com/boyapple/go-database/gorm/condition"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Client[T schema.Tabler] interface {
	// Get 获取单条数据
	Get(ctx context.Context, condition condition.Condition) (T, error)
	// List 获取列表数据
	List(ctx context.Context, opts ...Option) ([]T, error)
	// Count 获取总数
	Count(ctx context.Context, opts ...Option) (int64, error)
	// Create 创建
	Create(ctx context.Context, t T, opts ...Option) error
	// Update 更新
	Update(ctx context.Context, t T, condition condition.Condition) error

	UpdateKeyValue(ctx context.Context, keyValue map[string]interface{}, condition condition.Condition) error
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

func (i *impl[T]) Get(ctx context.Context, condition condition.Condition) (T, error) {
	var t T
	db, err := i.getDB(ctx)
	if err != nil {
		return t, err
	}
	scope, err := condition.Compile()
	if err != nil {
		return t, err
	}
	if err = db.Scopes(scope).First(&t).Error; err != nil {
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
	//opt := i.getOptions(opts...)
	db, err := i.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Create(t).Error
}

func (i *impl[T]) Update(ctx context.Context, t T, condition condition.Condition) error {
	db, err := i.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Updates(t).Error
}

func (i *impl[T]) UpdateKeyValue(ctx context.Context, keyValue map[string]interface{}, condition condition.Condition) error {
	db, err := i.getDB(ctx)
	if err != nil {
		return err
	}
	where, err := condition.Compile()
	if err != nil {
		return err
	}
	return db.Scopes(where).Updates(keyValue).Error
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
