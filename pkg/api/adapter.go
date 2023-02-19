package api

import (
	"auth/internal/app"
	"github.com/rs/zerolog"
	"net/http"
)

type adapter struct {
	logger     *zerolog.Logger
	config     *Config
	httpClient *http.Client
}

func NewAdapter(logger *zerolog.Logger, config *Config) app.Api {
	a := &adapter{
		logger:     logger,
		config:     config,
		httpClient: &http.Client{},
	}

	return a
}
