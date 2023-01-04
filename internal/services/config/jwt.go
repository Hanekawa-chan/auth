package config

type JWTConfig struct {
	SecretKey string `envconfig:"JWT_SECRET_KEY"`
}
