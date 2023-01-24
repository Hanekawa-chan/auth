package database

import "time"

type Config struct {
	Host     string `envconfig:"POSTGRES_HOST" required:"yes"`
	Port     int    `envconfig:"POSTGRES_PORT" required:"yes"`
	User     string `envconfig:"POSTGRES_USER" required:"yes"`
	Password string `envconfig:"POSTGRES_PASSWORD" required:"yes"`
	Name     string `envconfig:"POSTGRES_NAME" required:"yes"`

	MaxOpenConns    int           `envconfig:"POSTGRES_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int           `envconfig:"POSTGRES_MAX_IDLE_CONNS" default:"10"`
	ConnMaxLifeTime time.Duration `envconfig:"POSTGRES_CONN_MAX_LIFE_TIME" default:"5m"`
}
