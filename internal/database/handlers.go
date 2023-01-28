package database

import (
	"context"
	"github.com/google/uuid"
	"github.com/kanji-team/auth/internal/app"
	"github.com/kanji-team/auth/internal/database/models"
)

func (a *adapter) CreateUser(ctx context.Context, user *app.Credentials) error {
	var err error
	query := "insert into credentials (id, email, password, verified_email) values(:id, :email, :password, :verified_email)"

	_, err = a.db.NamedExecContext(ctx, query, models.CredentialsToDB(user))
	if err != nil {
		return err
	}
	return err
}

func (a *adapter) CreateGoogle(ctx context.Context, creds *app.Google) error {
	var err error
	query := "insert into google (id, email, google_id) values(:id, :email, :google_id)"

	_, err = a.db.NamedExecContext(ctx, query, models.GoogleToDB(creds))
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
