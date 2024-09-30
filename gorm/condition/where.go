package condition

import (
	"fmt"
	"gorm.io/gorm"
)

// ID id条件
type ID uint64

func (id ID) Compile() (func(*gorm.DB) *gorm.DB, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid id")
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}, nil
}

type Eq map[string]interface{}

func (eq Eq) Compile() (func(*gorm.DB) *gorm.DB, error) {
	return func(db *gorm.DB) *gorm.DB {
		for k, v := range eq {
			db.Where(fmt.Sprintf("%s = ?", k), v)
		}
		return db
	}, nil
}

func NewInCondition(field string, value []interface{}) Condition {
	return &In{
		Field: field,
		Value: value,
	}
}

type In struct {
	Field string
	Value []interface{}
}

func (in In) Compile() (func(*gorm.DB) *gorm.DB, error) {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IN (?)", in.Field), in.Value)
	}, nil
}
