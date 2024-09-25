package gorm

import (
	"github.com/boyapple/go-common/xmux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbMux = xmux.New[string, *gorm.DB]()

func New(name string, cfg *Config) (*gorm.DB, error) {
	db, err := dbMux.Get(name)
	if err == nil {
		return db, nil
	}
	db, err = gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	dbMux.Register(name, db)
	return db, nil
}
