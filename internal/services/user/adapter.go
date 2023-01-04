package user

import (
	"bytes"
	"context"
	"encoding/json"
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
	var user *models.User
	var err error

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := a.config.User.Address + "/api/v1/user/create"
	resp, err := a.client.Post(url, "text/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, err
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
