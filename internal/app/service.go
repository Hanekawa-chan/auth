package app

import (
	"github.com/rs/zerolog"
)

type service struct {
	logger       *zerolog.Logger
	jwtGenerator JWT
	db           Database
	user         User
	api          Api
}

func NewService(logger *zerolog.Logger, user User, api Api, generator JWT, database Database) Service {
	return &service{
		logger:       logger,
		db:           database,
		user:         user,
		jwtGenerator: generator,
		api:          api,
	}
}
