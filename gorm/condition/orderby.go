package condition

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewOrderByCondition(column string, desc bool) Condition {
	return &OrderBy{
		Column: column,
		Desc:   desc,
	}
}

// OrderBy 排序条件
type OrderBy struct {
	Column string
	Desc   bool
}

func (o *OrderBy) Compile() (func(*gorm.DB) *gorm.DB, error) {
	if o.Column == "" {
		return nil, fmt.Errorf("invalid order filed")
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(clause.OrderByColumn{
			Column: clause.Column{Name: o.Column},
			Desc:   o.Desc,
		})
	}, nil
}
