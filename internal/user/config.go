package user

import "time"

type Config struct {
	Address string        `envconfig:"USER_ADDRESS"`
	Timeout time.Duration `envconfig:"USER_TIMEOUT"`
}
