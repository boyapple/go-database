package porm

import (
	"github.com/boyapple/go-database/porm/filter"
	"github.com/boyapple/go-database/porm/pb"
)

func NewOptions() *Options {
	return &Options{
		TimeFieldFilter: filter.NewEmptyFieldFilter(),
	}
}

// Options 选项
type Options struct {
	Fields          []string // 自定义select,update,insert字段
	Table           string   // 表名
	Where           string   // where条件
	Args            []interface{}
	Page            *pb.Page      // 分页
	OrderBy         []*pb.OrderBy // 排序
	Join            string
	TimeFieldFilter filter.Filter // 时间字段过滤器
}

// Option 选项
type Option func(options *Options)

func WithFields(fields []string) Option {
	return func(o *Options) {
		o.Fields = fields
	}
}

func WithTable(table string) Option {
	return func(o *Options) {
		o.Table = table
	}
}

func WithPage(page *pb.Page) Option {
	return func(o *Options) {
		o.Page = page
	}
}

func WithOrderBy(orderBy ...*pb.OrderBy) Option {
	return func(o *Options) {
		o.OrderBy = orderBy
	}
}

func WithWhereArgs(where string, args ...interface{}) Option {
	return func(o *Options) {
		o.Where = where
		o.Args = args
	}
}

func WithJoin(join string) Option {
	return func(o *Options) {
		o.Join = join
	}
}

func WithTimeField(timeField []string) Option {
	return func(os *Options) {
		os.TimeFieldFilter = filter.NewTimeFieldFilter(timeField)
	}
}

func WithTimeFieldFilter(filter filter.Filter) Option {
	return func(o *Options) {
		o.TimeFieldFilter = filter
	}
}
