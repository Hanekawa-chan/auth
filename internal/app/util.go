package app

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"time"
	"unicode/utf8"
)

const Cost = 12
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

func (a *service) validatePair(login string, password string) error {
	err := a.validateEmail(login)
	if err != nil {
		return ErrValidation
	}

	err = a.validatePassword(password)
	if err != nil {
		return err
	}

	return nil
}

func (a *service) hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), Cost)
}

func (a *service) generateAccessToken(userID uuid.UUID) (string, error) {
	claims := make(map[string]interface{})
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour).Unix()

	token, err := a.jwtGenerator.Generate(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (a *service) generateRefreshToken(userID uuid.UUID) (string, error) {
	claims := make(map[string]interface{})
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(168 * time.Hour).Unix()

	token, err := a.jwtGenerator.Generate(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}
