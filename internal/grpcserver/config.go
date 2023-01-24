package grpcserver

import "time"

type Config struct {
	Address           string        `envconfig:"GRPC_ADDRESS"`
	SecretKey         string        `envconfig:"JWT_SECRET_KEY"`
	MaxConnectionIdle time.Duration `envconfig:"MAX_CONNECTION_IDLE"`
	Timeout           time.Duration `envconfig:"TIMEOUT"`
	MaxConnectionAge  time.Duration `envconfig:"MAX_CONNECTION_AGE"`
}
