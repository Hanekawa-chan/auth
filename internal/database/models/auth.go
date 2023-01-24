package models

import (
	"github.com/google/uuid"
	"github.com/kanji-team/auth/internal/app"
)

type Credentials struct {
	Id       string
	Login    string
	Password string
	AuthHash string `db:"auth_hash"`
}

func (c *Credentials) ToDomain() (*app.Credentials, error) {
	id, err := uuid.Parse(c.Id)
	if err != nil {
		return nil, err
	}

	return &app.Credentials{
		Id:       id,
		Email:    c.Login,
		Password: c.Password,
		AuthHash: c.AuthHash,
	}, nil
}

type Google struct {
	Id       string
	Email    string
	GoogleId string `db:"google_id"`
}

func (g *Google) ToDomain() (*app.Google, error) {
	id, err := uuid.Parse(g.Id)
	if err != nil {
		return nil, err
	}

	return &app.Google{
		Id:       id,
		Email:    g.Email,
		GoogleId: g.GoogleId,
	}, nil
}
