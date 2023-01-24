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
	Login    string
	Password string
	AuthHash string `db:"auth_hash"`
}

type Google struct {
	Id       uuid.UUID
	Email    string
	GoogleId string `db:"google_id"`
}
