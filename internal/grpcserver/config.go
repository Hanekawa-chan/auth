package grpcserver

type Config struct {
	Address   string `envconfig:"GRPC_ADDRESS"`
	SecretKey string `envconfig:"JWT_SECRET_KEY"`
}
