package user

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Hanekawa-chan/kanji-auth/internal/app"
	"github.com/Hanekawa-chan/kanji-auth/internal/services/models"
	"github.com/rs/zerolog"
	"net/http"
)

type adapter struct {
	logger *zerolog.Logger
	config *app.Config
	client *http.Client
}

func (a *adapter) CreateUser(ctx context.Context, req *models.CreateUserRequest) (string, error) {
	var user models.CreateUserResponse
	var err error

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	url := a.config.User.Address + "/api/v1/user/create"
	resp, err := a.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return "", err
	}

	return user.UserId, err
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
