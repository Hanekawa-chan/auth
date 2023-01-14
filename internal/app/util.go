package app

import (
	"crypto/rand"
	"errors"
	"github.com/google/uuid"
	"net/mail"
	"unicode/utf8"
)

const MinSymbols = 8
const MaxSymbols = 32

func (a *service) validatePassword(password string) error {
	length := utf8.RuneCountInString(password)
	if length < MinSymbols || length > MaxSymbols {
		return errors.New("invalid password length")
	}
	return nil
}

func (a *service) validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

func (a *service) generateJWT(userID uuid.UUID) (string, error) {
	claims := make(map[string]interface{})
	claims["user_id"] = userID

	token, err := a.jwtGenerator.Generate(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

func generateAuthHash() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	hash := string(b)

	return hash, err
}
