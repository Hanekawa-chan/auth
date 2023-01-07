package app

import (
	"github.com/Hanekawa-chan/kanji-auth/internal/services/config"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Logger     *LoggerConfig
	HTTPServer *config.HTTPConfig
	Auth       *config.AuthConfig
	DB         *config.DBConfig
	JWTConfig  *config.JWTConfig
	User       *config.UserConfig
}

type LoggerConfig struct {
	LogLevel string `default:"debug"`
}

func Parse() (*Config, error) {
	cfg := Config{}
	logger := LoggerConfig{}
	db := config.DBConfig{}
	auth := config.AuthConfig{}
	httpConfig := config.HTTPConfig{}
	jwtConfig := config.JWTConfig{}
	userConfig := config.UserConfig{}

	err := envconfig.Process("kanji_auth", &logger)
	if err != nil {
		log.Err(err).Msg("logger config error")
		return nil, err
	}

	err = envconfig.Process("kanji_auth", &db)
	if err != nil {
		log.Err(err).Msg("db config error")
		return nil, err
	}

	err = envconfig.Process("kanji_auth", &auth)
	if err != nil {
		log.Err(err).Msg("auth config error")
		return nil, err
	}

	err = envconfig.Process("kanji_auth", &httpConfig)
	if err != nil {
		log.Err(err).Msg("http config error")
		return nil, err
	}

	err = envconfig.Process("kanji_auth", &jwtConfig)
	if err != nil {
		log.Err(err).Msg("http config error")
		return nil, err
	}

	err = envconfig.Process("kanji_auth", &userConfig)
	if err != nil {
		log.Err(err).Msg("http config error")
		return nil, err
	}

	cfg.Auth = &auth
	cfg.DB = &db
	cfg.Logger = &logger
	cfg.HTTPServer = &httpConfig
	cfg.JWTConfig = &jwtConfig
	cfg.User = &userConfig

	return &cfg, nil
}
