package db

import (
	"fmt"
	"gorm.io/gorm"
)

type Conditions []Condition

type Condition interface {
	Where() (func(db *gorm.DB) *gorm.DB, error)
}

func (c Conditions) Build() ([]func(db *gorm.DB) *gorm.DB, error) {
	scopes := make([]func(db *gorm.DB) *gorm.DB, 0, len(c))
	for _, w := range c {
		where, err := w.Where()
		if err != nil {
			return nil, err
		}
		scopes = append(scopes, where)
	}
	return scopes, nil
}

type Eq map[string]any

func (eq Eq) Where() (func(db *gorm.DB) *gorm.DB, error) {
	return func(db *gorm.DB) *gorm.DB {
		for k, v := range eq {
			db.Where(fmt.Sprintf("%s = ?", k), v)
		}
		return db
	}, nil
}
