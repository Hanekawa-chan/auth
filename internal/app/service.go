package app

import (
	"context"
	"github.com/rs/zerolog"
	"kanji-auth/internal/services/models"
)

type service struct {
	logger *zerolog.Logger
	cfg    *Config
	auth   Auth
}

func NewService(logger *zerolog.Logger, cfg *Config, auth Auth) Service {
	return service{
		logger: logger,
		cfg:    cfg,
		auth:   auth,
	}
}

func (s service) Auth(ctx context.Context, req *models.AuthRequest) (*models.Session, error) {
	return s.auth.Auth(ctx, req)
}

func (s service) Signup(ctx context.Context, req *models.SignupRequest) (*models.Session, error) {
	return s.auth.Signup(ctx, req)
}
