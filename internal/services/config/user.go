package config

import "time"

type UserConfig struct {
	Address string        `envconfig:"USER_ADDRESS"`
	Timeout time.Duration `envconfig:"USER_TIMEOUT"`
}
