package db

import (
	"github.com/boyapple/go-common/xmux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbMux = xmux.New[string, *gorm.DB]()

type Config struct {
	Dsn string
}

func Register(name string, cfg *Config) error {
	db, err := gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	dbMux.Register(name, db.Debug())
	return nil
}
