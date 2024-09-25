package gorm

type Config struct {
	Dsn string
}

type ConfigOption func(*Config)
