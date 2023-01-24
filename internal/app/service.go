package app

import (
	"github.com/kanji-team/auth/internal/app/config"
	jwt "github.com/kanji-team/jwt"
	"github.com/rs/zerolog"
)

type service struct {
	logger       *zerolog.Logger
	config       *config.Config
	db           Database
	jwtGenerator *jwt.Generator
	user         User
}

func NewService(logger *zerolog.Logger, cfg *config.Config, user User, generator *jwt.Generator, database Database) Service {
	return &service{
		logger:       logger,
		config:       cfg,
		db:           database,
		user:         user,
		jwtGenerator: generator,
	}
}
