package httpserver

import (
	"context"
	"encoding/json"
	"github.com/Hanekawa-chan/kanji-auth/internal/services/models"
	"net/http"
)

func (a *adapter) Auth(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	req := models.AuthRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = a.service.Auth(ctx, &req)
	if err != nil {
		return
	}
}

func (a *adapter) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	req := models.SignupRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = a.service.Signup(ctx, &req)
	if err != nil {
		return
	}
}

func (a *adapter) Link(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	req := models.AuthRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.service.Link(ctx, &req)
	if err != nil {
		return
	}
}
