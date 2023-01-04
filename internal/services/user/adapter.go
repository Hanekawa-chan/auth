package user

import (
	"context"
	"github.com/rs/zerolog"
	"kanji-auth/internal/app"
	"kanji-auth/internal/services/models"
	"net/http"
)

type adapter struct {
	logger *zerolog.Logger
	config *app.Config
	client *http.Client
}

func (a adapter) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserClient(logger *zerolog.Logger, config *app.Config) app.User {
	client := &http.Client{
		Timeout: config.User.Timeout,
	}

	return &adapter{
		logger: logger,
		config: config,
		client: client,
	}
}
