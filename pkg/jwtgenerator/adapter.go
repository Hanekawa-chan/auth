package jwtgenerator

import (
	"auth/internal/app"
	"github.com/Hanekawa-chan/jwt"
	jwtx "github.com/golang-jwt/jwt/v4"
)

type adapter struct {
	config    *Config
	generator *jwt.Generator
}

func NewAdapter(config *Config) (app.JWT, error) {
	a := &adapter{
		config: config,
	}

	generator, err := jwt.New(config.SecretKey)
	if err != nil {
		return nil, err
	}

	a.generator = generator

	return a, nil
}

func (a adapter) Generate(claims map[string]interface{}) (string, error) {
	return a.generator.Generate(claims)
}

func (a adapter) ParseToken(token string) (jwtx.MapClaims, error) {
	return a.generator.ParseToken(token)
}
