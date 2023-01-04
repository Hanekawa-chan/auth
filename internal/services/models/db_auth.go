package models

import "github.com/google/uuid"

type Credentials struct {
	ID       uuid.UUID `db:"id"`
	Email    string
	Password string
	AuthHash string `db:"auth_hash"`
}

type Google struct {
	ID       uuid.UUID `db:"id"`
	Email    string
	GoogleID string `db:"google_id"`
}
