package app

import (
	"context"
	"github.com/Hanekawa-chan/kanji-auth/internal/services/models"
	"github.com/google/uuid"
)

type Service interface {
	Auth(ctx context.Context, req *models.AuthRequest) (*models.Session, error)
	Signup(ctx context.Context, req *models.SignupRequest) (*models.Session, error)
	Link(ctx context.Context, req *models.AuthRequest) error
}

type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type Database interface {
	GetUserByAuthHash(ctx context.Context, hash string) (*models.Credentials, error)
	UpdateId(ctx context.Context, id uuid.UUID, hash string) error
	RemoveAuthHash(ctx context.Context, id uuid.UUID) error
	CreateUser(ctx context.Context, user *models.Credentials) error
	GetUserByEmail(ctx context.Context, login string) (*models.Credentials, error)
	GetUserByGoogleEmail(ctx context.Context, email string) (*models.Credentials, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.Credentials, error)
	CreateGoogle(ctx context.Context, creds *models.Google) error
}

type JWTGenerator interface {
	Generate(claims map[string]interface{}) (string, error)
}

type User interface {
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
}
