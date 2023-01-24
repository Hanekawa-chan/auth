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
	Id       uuid.UUID
	Email    string
	Password string
	AuthHash string
}

type Google struct {
	Id       uuid.UUID
	Email    string
	GoogleId string
}
