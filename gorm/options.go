package gorm

import "gorm.io/gorm"

type Options struct {
	Page      *Page
	Condition Condition
}

type Option func(*Options)

func WithCondition(condition Condition) Option {
	return func(o *Options) {
		o.Condition = condition
	}
}

func WithPage(offset, limit int) Option {
	return func(o *Options) {
		o.Page = &Page{
			Offset: offset,
			Limit:  limit,
		}
	}
}

type Condition interface {
	Where() func(db *gorm.DB) *gorm.DB
}

type Page struct {
	Offset int
	Limit  int
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
