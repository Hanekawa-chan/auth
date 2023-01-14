package httpserver

import (
	"context"
	"github.com/Hanekawa-chan/kanji-auth/internal/services/models"
	"net/http"
)

func (a *adapter) Auth(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	req := models.AuthRequest{}
	_, err := a.service.Auth(ctx, &req)
	if err != nil {
		return
	}
}

func (a *adapter) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	req := models.SignupRequest{}
	_, err := a.service.Signup(ctx, &req)
	if err != nil {
		return
	}
}

func (a *adapter) Link(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	req := models.AuthRequest{}
	err := a.service.Link(ctx, &req)
	if err != nil {
		return
	}
}
