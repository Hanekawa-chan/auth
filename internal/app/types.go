package app

import "github.com/google/uuid"

type GoogleAuthUser struct {
	ID            string
	Email         string
	Name          string
	GivenName     string
	FamilyName    string
	Picture       string
	Locale        string
	VerifiedEmail bool
}

type Credentials struct {
	Id            uuid.UUID
	Email         string
	Password      string
	VerifiedEmail bool
	IssuedAt      int64
}

type Google struct {
	Id       uuid.UUID
	Email    string
	GoogleId string
}

type SignUpRequest struct {
	Name  string
	Email string
}
