package user

import (
	"context"
	"github.com/Hanekawa-chan/kanji-auth/proto/services"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type adapter struct {
	services.InternalUserServiceClient
	logger *zerolog.Logger
	config *Config
}

func (a *adapter) CreateUser(ctx context.Context, req *services.CreateUserRequest) (*services.CreateUserResponse, error) {

}

func NewUserClient(logger *zerolog.Logger, config *Config) domain.User {
	conn, err := grpc.Dial(config.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("create user service client")
	}

	client := services.NewInternalUserServiceClient(conn)

	return &adapter{
		logger:                    logger,
		config:                    config,
		InternalUserServiceClient: client,
	}
}
