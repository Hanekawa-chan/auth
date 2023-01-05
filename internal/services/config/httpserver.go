package config

type HTTPConfig struct {
	Address   string `envconfig:"HTTP_SERVER_ADDRESS" default:":6080"`
	RateLimit int    `envconfig:"HTTP_RATE_LIMIT" default:"20"`
}
