package porm

import (
	"context"
	"database/sql"

	"github.com/golang/protobuf/proto"
)

type (
	// NextFunc 行解析拦截器
	NextFunc func(rows *sql.Rows) error
)

// Client 执行接口
type Client interface {
	// First 查询一行,必须指定条件
	First(ctx context.Context, message proto.Message, opts ...Option) error
	// List 查询列表
	List(ctx context.Context, dst interface{}, opts ...Option) error
	// Insert 插入一条数据
	Insert(ctx context.Context, message proto.Message, opts ...Option) (int64, error)
	// Update 更新一条数据,必须指定条件
	Update(ctx context.Context, message proto.Message, opts ...Option) (int64, error)
	// Exec 自定义sql执行,自行保证sql注入
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	// Query 自定义查询
	Query(ctx context.Context, query string, next NextFunc, args ...any) error
}
