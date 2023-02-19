package jwtgenerator

type Config struct {
	SecretKey string `envconfig:"SECRET_KEY"`
}
