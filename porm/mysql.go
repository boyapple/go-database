package porm

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
)

func New(dsn string) (Client, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &mysqlCli{db: db}, nil
}

type mysqlCli struct {
	db *sql.DB
}

func (m *mysqlCli) First(ctx context.Context, message proto.Message, opts ...Option) error {
	options := NewOptions()
	for _, o := range opts {
		o(options)
	}
	query, args, err := NewSelectBuilder(SelectOne).Build(message, options)
	if err != nil {
		return err
	}
	next := func(rows *sql.Rows) error {
		return ParseRowsProto(rows, message, options.TimeFieldFilter)
	}
	return m.Query(ctx, query, next, args...)
}

func (m *mysqlCli) List(ctx context.Context, dst interface{}, opts ...Option) error {
	options := NewOptions()
	for _, o := range opts {
		o(options)
	}
	value := reflect.ValueOf(dst)
	// 判断是否为指针
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return fmt.Errorf("dst is not pointer")
	}
	direct := reflect.Indirect(value)
	slice, err := valueType(value.Type(), reflect.Slice)
	if err != nil {
		return err
	}
	baseType := deref(slice.Elem())
	message, ok := reflect.New(baseType).Interface().(proto.Message)
	if !ok {
		return fmt.Errorf("struct not proto.message")
	}
	query, args, err := NewSelectBuilder(SelectList).Build(message, options)
	if err != nil {
		return err
	}
	next := func(rows *sql.Rows) error {
		val := reflect.New(baseType)
		okMessage, ook := val.Interface().(proto.Message)
		if !ook {
			return fmt.Errorf("struct not proto.message")
		}
		if err = ParseRowsProto(rows, okMessage, options.TimeFieldFilter); err != nil {
			return err
		}
		direct.Set(reflect.Append(direct, val))
		return nil
	}
	if err = m.Query(ctx, query, next, args...); err != nil {
		return err
	}
	return m.Count(ctx, options)
}

func (m *mysqlCli) Insert(ctx context.Context, message proto.Message, opts ...Option) (int64, error) {
	options := NewOptions()
	for _, o := range opts {
		o(options)
	}
	query, args, err := NewInsertBuilder().Build(message, options)
	if err != nil {
		return 0, err
	}
	result, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (m *mysqlCli) Update(ctx context.Context, message proto.Message, opts ...Option) (int64, error) {
	options := NewOptions()
	for _, o := range opts {
		o(options)
	}
	query, args, err := NewUpdateBuilder().Build(message, options)
	if err != nil {
		return 0, err
	}
	result, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *mysqlCli) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return m.db.ExecContext(ctx, query, args...)
}

func (m *mysqlCli) Query(ctx context.Context, query string, next NextFunc, args ...any) error {
	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err = next(rows); err != nil {
			return err
		}
	}
	return nil
}

func (m *mysqlCli) Count(ctx context.Context, opts *Options) error {
	if opts.Page != nil && opts.Page.Offset > 0 && opts.Page.Limit > 0 {
		query, args, err := NewSelectBuilder(SelectCount).Build(nil, opts)
		if err != nil {
			return err
		}
		return m.db.QueryRowContext(ctx, query, args...).Scan(&opts.Page.Total)
	}
	return nil
}
