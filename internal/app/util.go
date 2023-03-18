package app

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"time"
	"unicode/utf8"
)

const Cost = 12
const MinSymbols = 10
const MaxSymbols = 32
const AccessExp = 30 * time.Minute
const RefreshExp = 2 * 31 * 24 * time.Hour

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

func (a *service) generateTokens(ctx context.Context, userID uuid.UUID) (string, string, error) {
	issuedAt := time.Now()

	accessToken, err := a.generateAccessToken(userID, issuedAt)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := a.generateRefreshToken(issuedAt)
	if err != nil {
		return "", "", err
	}

	err = a.db.UpdateIssuedAt(ctx, userID, issuedAt.Unix())
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (a *service) generateAccessToken(userID uuid.UUID, issuedAt time.Time) (string, error) {
	claims := make(map[string]interface{})
	claims["user_id"] = userID
	claims["iat"] = issuedAt.Unix()
	claims["exp"] = issuedAt.Add(AccessExp).Unix()

	token, err := a.jwtGenerator.Generate(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (a *service) generateRefreshToken(issuedAt time.Time) (string, error) {
	claims := make(map[string]interface{})
	claims["iat"] = issuedAt.Unix()
	claims["exp"] = issuedAt.Add(RefreshExp).Unix()

	token, err := a.jwtGenerator.Generate(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (a *service) parseAccessToken(token string) (uuid.UUID, int64, error) {
	claims, err := a.jwtGenerator.ParseToken(token)
	if err != nil {
		if err == jwt.ErrTokenExpired {
			tokenStruct, _, err := (&jwt.Parser{}).ParseUnverified(token, jwt.MapClaims{})
			if err != nil {
				return uuid.UUID{}, 0, err
			}

			claims = tokenStruct.Claims.(jwt.MapClaims)
		} else {
			return uuid.UUID{}, 0, err
		}
	}

	userID := uuid.UUID{}
	if id, ok := claims["user_id"].(uuid.UUID); !ok {
		return uuid.UUID{}, 0, ErrType
	} else {
		userID = id
	}

	var issuedAt int64 = 0
	if iat, ok := claims["iat"].(int64); !ok {
		return uuid.UUID{}, 0, ErrType
	} else {
		issuedAt = iat
	}

	return userID, issuedAt, err
}

func (a *service) parseRefreshToken(token string) error {
	_, err := a.jwtGenerator.ParseToken(token)
	if err != nil {
		return err
	}

	return nil
}
