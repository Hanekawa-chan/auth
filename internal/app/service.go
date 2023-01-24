package app

import (
	"github.com/Hanekawa-chan/kanji-auth/internal/app/config"
	"github.com/Hanekawa-chan/kanji-auth/proto/services"
	"github.com/rs/zerolog"
)

type service struct {
	logger       *zerolog.Logger
	config       *config.Config
	db           Database
	jwtGenerator JWTGenerator
	user         services.InternalUserServiceClient
}

func NewService(logger *zerolog.Logger, cfg *config.Config, user services.InternalUserServiceClient, generator JWTGenerator, database Database) Service {
	return &service{
		logger:       logger,
		config:       cfg,
		db:           database,
		user:         user,
		jwtGenerator: generator,
	}
}
