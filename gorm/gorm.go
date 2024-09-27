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

type ConfigOption func(*Config)

func Get(name string) (*gorm.DB, error) {
	return dbMux.Get(name)
}

func Register(name string, cfg *Config) error {
	db, err := New(cfg.Dsn)
	if err != nil {
		return err
	}
	dbMux.Register(name, db)
	return nil
}

func New(dsn string, cfgOpts ...ConfigOption) (*gorm.DB, error) {
	cfg := &Config{}
	for _, o := range cfgOpts {
		o(cfg)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
