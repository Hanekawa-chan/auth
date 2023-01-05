package database

import (
	"context"
	"github.com/google/uuid"
	"kanji-auth/internal/services/models"
)

func (a *adapter) UpdateId(ctx context.Context, id uuid.UUID, hash string) error {
	var err error
	query := "update credentials set id=$1 where auth_hash=$2"

	_, err = a.db.ExecContext(ctx, query, id, hash)
	if err != nil {
		return err
	}
	return err
}

func (a *adapter) RemoveAuthHash(ctx context.Context, id uuid.UUID) error {
	var err error
	query := "update credentials set auth_hash=null where id=$1"

	_, err = a.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return err
}

func (a *adapter) CreateUser(ctx context.Context, user *models.Credentials) error {
	var err error
	query := "insert into credentials (id, email, password, auth_hash) values(:id, :email, :password, :auth_hash)"

	_, err = a.db.NamedExecContext(ctx, query, &user)
	if err != nil {
		return err
	}
	return err
}

func (a *adapter) GetUserByEmail(ctx context.Context, login string) (*models.Credentials, error) {
	var err error
	creds := models.Credentials{}
	query := "select * from credentials where email=$1"

	err = a.db.SelectContext(ctx, &creds, query, login)
	if err != nil {
		return nil, err
	}
	return &creds, err
}

func (a *adapter) GetUserByGoogleEmail(ctx context.Context, email string) (*models.Credentials, error) {
	var err error
	creds := models.Credentials{}
	query := "select * from credentials where email=$1"

	err = a.db.SelectContext(ctx, &creds, query, email)
	if err != nil {
		return nil, err
	}
	return &creds, err
}

func (a *adapter) GetUserByID(ctx context.Context, id uuid.UUID) (*models.Credentials, error) {
	var err error
	creds := models.Credentials{}
	query := "select * from credentials where id=$1"

	err = a.db.SelectContext(ctx, &creds, query, id)
	if err != nil {
		return nil, err
	}
	return &creds, err
}

func (a *adapter) GetUserByAuthHash(ctx context.Context, hash string) (*models.Credentials, error) {
	var err error
	creds := models.Credentials{}
	query := "select * from credentials where auth_hash=$1"

	err = a.db.SelectContext(ctx, &creds, query, hash)
	if err != nil {
		return nil, err
	}
	return &creds, err
}
