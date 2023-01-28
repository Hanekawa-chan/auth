package api

import (
	"github.com/kanji-team/auth/internal/app"
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
