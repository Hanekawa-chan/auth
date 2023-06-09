package models

import (
	"auth/internal/app"
	"database/sql"
	"github.com/google/uuid"
)

type Credentials struct {
	Id            uuid.UUID      `db:"id"`
	Email         string         `db:"email"`
	Password      sql.NullString `db:"password"`
	VerifiedEmail bool           `db:"verified_email"`
	IssuedAt      int64          `db:"iat"`
}

func (c *Credentials) ToDomain() *app.Credentials {
	password := ""

	if c.Password.Valid {
		password = c.Password.String
	}

	return &app.Credentials{
		Id:            c.Id,
		Email:         c.Email,
		Password:      password,
		VerifiedEmail: c.VerifiedEmail,
	}
}

func FromCredentials(c *app.Credentials) *Credentials {
	password := sql.NullString{
		String: "",
		Valid:  false,
	}

	if c.Password != "" {
		password.Valid = true
		password.String = c.Password
	}

	return &Credentials{
		Id:            c.Id,
		Email:         c.Email,
		Password:      password,
		VerifiedEmail: c.VerifiedEmail,
	}
}

type Google struct {
	Id       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	GoogleId string    `db:"google_id"`
}

func (g *Google) ToDomain() *app.Google {
	return &app.Google{
		Id:       g.Id,
		Email:    g.Email,
		GoogleId: g.GoogleId,
	}
}

func FromGoogle(g *app.Google) *Google {
	return &Google{
		Id:       g.Id,
		Email:    g.Email,
		GoogleId: g.GoogleId,
	}
}
