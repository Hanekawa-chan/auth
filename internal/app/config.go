package app

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"kanji-auth/internal/services/config"
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
	var cfg *Config
	var logger *LoggerConfig
	var db *config.DBConfig
	var auth *config.AuthConfig
	var httpConfig *config.HTTPConfig
	var jwtConfig *config.JWTConfig
	var userConfig *config.UserConfig

	err := envconfig.Process("kanji_auth", logger)
	if err != nil {
		log.Err(err).Msg("logger config error")
		return nil, err
	}

	err = envconfig.Process("kanji_auth", db)
	if err != nil {
		log.Err(err).Msg("db config error")
		return nil, err
	}

	err = envconfig.Process("kanji_auth", auth)
	if err != nil {
		log.Err(err).Msg("auth config error")
		return nil, err
	}

	err = envconfig.Process("kanji_auth", httpConfig)
	if err != nil {
		log.Err(err).Msg("http config error")
		return nil, err
	}

	err = envconfig.Process("kanji_auth", jwtConfig)
	if err != nil {
		log.Err(err).Msg("http config error")
		return nil, err
	}

	err = envconfig.Process("kanji_auth", userConfig)
	if err != nil {
		log.Err(err).Msg("http config error")
		return nil, err
	}

	cfg.Auth = auth
	cfg.DB = db
	cfg.Logger = logger
	cfg.HTTPServer = httpConfig
	cfg.JWTConfig = jwtConfig
	cfg.User = userConfig

	return cfg, nil
}
