package config

type HTTPConfig struct {
	Address   string `envconfig:"HTTP_SERVER_ADDRESS" envDefault:":6000"`
	RateLimit int    `envconfig:"HTTP_RATE_LIMIT" envDefault:"20"`
}
