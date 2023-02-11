package api

import (
	"auth/internal/app"
	"github.com/rs/zerolog"
)

type adapter struct {
	logger *zerolog.Logger
	config *Config
}

func NewAdapter(logger *zerolog.Logger, config *Config) app.Api {
	a := &adapter{
		logger: logger,
		config: config,
	}

	return a
}
