package config

import "time"

type DBConfig struct {
	Host     string `envconfig:"POSTGRES_HOST" required:"yes"`
	Port     int    `envconfig:"POSTGRES_PORT" required:"yes"`
	User     string `envconfig:"POSTGRES_USER" required:"yes"`
	Password string `envconfig:"POSTGRES_PASSWORD" required:"yes"`
	Name     string `envconfig:"POSTGRES_NAME" required:"yes"`

	MaxOpenConns    int           `envconfig:"POSTGRES_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `envconfig:"POSTGRES_MAX_IDLE_CONNS" envDefault:"10"`
	ConnMaxLifeTime time.Duration `envconfig:"POSTGRES_CONN_MAX_LIFE_TIME" envDefault:"5m"`

	//MigrationsSourceURL string `env:"POSTGRES_MIGRATIONS_SOURCE_URL" envDefault:"file://migrations"`
}
