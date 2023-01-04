package app

import (
	"context"
	"github.com/google/uuid"
	"kanji-auth/internal/services/models"
)

type Service interface {
	Auth(ctx context.Context, req *models.AuthRequest) (*models.Session, error)
	Signup(ctx context.Context, req *models.SignupRequest) (*models.Session, error)
}

type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type Auth interface {
	Auth(ctx context.Context, req *models.AuthRequest) (*models.Session, error)
	Signup(ctx context.Context, req *models.SignupRequest) (*models.Session, error)
}

type Database interface {
	GetUserByAuthHash(ctx context.Context, hash string) (*models.Credentials, error)
	UpdateId(ctx context.Context, id uuid.UUID, hash string) error
	RemoveAuthHash(ctx context.Context, id uuid.UUID) error
}

type JWTGenerator interface {
	Generate(claims map[string]interface{}) (string, error)
}

type User interface {
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
}
