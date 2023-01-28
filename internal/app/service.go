package app

import (
	jwt "github.com/kanji-team/jwt"
	"github.com/rs/zerolog"
)

type service struct {
	logger       *zerolog.Logger
	config       *Config
	jwtGenerator *jwt.Generator
	db           Database
	user         User
	api          Api
}

func NewService(logger *zerolog.Logger, cfg *Config, user User, api Api, generator *jwt.Generator, database Database) Service {
	return &service{
		logger:       logger,
		config:       cfg,
		db:           database,
		user:         user,
		jwtGenerator: generator,
		api:          api,
	}
}
