package models

import "github.com/google/uuid"

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
