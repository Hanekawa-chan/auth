package auth

import (
	"github.com/Hanekawa-chan/kanji-auth/internal/app"
	"github.com/rs/zerolog"
)

type adapter struct {
	logger       *zerolog.Logger
	config       *app.Config
	db           app.Database
	jwtGenerator app.JWTGenerator
	user         app.User
}

func NewAuth(logger *zerolog.Logger, db app.Database, jwtGenerator app.JWTGenerator, user app.User, config *app.Config) app.Auth {
	return &adapter{
		logger:       logger,
		db:           db,
		jwtGenerator: jwtGenerator,
		user:         user,
		config:       config,
	}
}
