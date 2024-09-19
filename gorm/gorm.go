package gorm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(name string, opts ...Option) (*gorm.DB, error) {
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	db, err := gorm.Open(mysql.Open(opt.Dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
