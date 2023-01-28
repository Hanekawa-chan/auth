package app

type Config struct {
	JWTSecretKey string `envconfig:"JWT_SECRET_KEY"`
}
