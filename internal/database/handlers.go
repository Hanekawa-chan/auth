package database

import (
	"auth/internal/app"
	"auth/internal/database/models"
	"context"
	"github.com/google/uuid"
)

func (a *adapter) CreateUser(ctx context.Context, user *app.Credentials) error {
	var err error
	query := "insert into credentials (id, email, password, verified_email) values($1, $2, $3, $4)"
	dbUser := models.FromCredentials(user)

	_, err = a.db.ExecContext(ctx, query, dbUser.Id, dbUser.Email, dbUser.Password, dbUser.VerifiedEmail)
	if err != nil {
		return err
	}
	return err
}

func (a *adapter) CreateGoogle(ctx context.Context, creds *app.Google) error {
	var err error
	query := "insert into google (id, email, google_id) values($1, $2, $3)"
	dbGoogle := models.FromGoogle(creds)

	_, err = a.db.ExecContext(ctx, query, dbGoogle.Id, dbGoogle.Email, dbGoogle.GoogleId)
	if err != nil {
		return err
	}
	return err
}

func (a *adapter) GetUserByEmail(ctx context.Context, login string) (*app.Credentials, error) {
	var err error
	creds := &models.Credentials{}
	query := "select * from credentials where email=$1"

	err = a.db.SelectContext(ctx, creds, query, login)
	if err != nil {
		return nil, err
	}
	return creds.ToDomain(), nil
}

func (a *adapter) GetUserByGoogleEmail(ctx context.Context, email string) (*app.Google, error) {
	var err error
	creds := &models.Google{}
	query := "select * from google where email=$1"

	err = a.db.SelectContext(ctx, creds, query, email)
	if err != nil {
		return nil, err
	}
	return creds.ToDomain(), nil
}

func (a *adapter) GetUserByID(ctx context.Context, id uuid.UUID) (*app.Credentials, error) {
	var err error
	creds := &models.Credentials{}
	query := "select * from credentials where id=$1"

	err = a.db.SelectContext(ctx, creds, query, id)
	if err != nil {
		return nil, err
	}
	return creds.ToDomain(), nil
}
