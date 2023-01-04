package database

import (
	"context"
	"github.com/google/uuid"
	"kanji-auth/internal/services/models"
)

func (a adapter) GetUserByAuthHash(ctx context.Context, hash string) (*models.Credentials, error) {
	//TODO implement me
	panic("implement me")
}

func (a adapter) UpdateId(ctx context.Context, id uuid.UUID, hash string) error {
	//TODO implement me
	panic("implement me")
}

func (a adapter) RemoveAuthHash(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
