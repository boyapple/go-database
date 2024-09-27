package db

type Options struct {
	Page         *Page
	Condition    Conditions
	OnlyColumn   []string // 用于on duplicate key update语句的条件字段
	UpdateColumn []string // 指定需要更新的列
}

type Page struct {
	Offset int   `json:"offset"`
	Limit  int   `json:"limit"`
	Count  int64 `json:"count"`
}

func (p *Page) GetOffset() int {
	if p.Offset <= 0 {
		return 1
	}
	return p.Offset
}

func (p *Page) GetLimit() int {
	if p.Limit <= 0 {
		return 10
	}
	return p.Limit
}

type Option func(*Options)

func WithCondition(condition ...Condition) Option {
	return func(o *Options) {
		o.Condition = condition
	}
}

func WithPage(page *Page) Option {
	return func(o *Options) {
		o.Page = page
	}
}
