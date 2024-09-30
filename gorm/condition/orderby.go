package condition

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewOrderByCondition(field string, desc bool) Condition {
	return &OrderBy{
		Field: field,
		Desc:  desc,
	}
}

// OrderBy 排序条件
type OrderBy struct {
	Field string
	Desc  bool
}

func (o *OrderBy) Compile() (func(*gorm.DB) *gorm.DB, error) {
	if o.Field == "" {
		return nil, fmt.Errorf("invalid order filed")
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(clause.OrderByColumn{
			Column: clause.Column{Name: o.Field},
			Desc:   o.Desc,
		})
	}, nil
}
