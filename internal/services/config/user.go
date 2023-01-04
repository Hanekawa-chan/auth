package config

import "time"

type UserConfig struct {
	Timeout time.Duration `envconfig:"USER_TIMEOUT"`
}
