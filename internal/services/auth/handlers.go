package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"kanji-auth/internal/services/models"
	"net/http"
)

func (a adapter) GetUserInfoFromGoogleAPI(ctx context.Context, code string) (*models.GoogleAuthUser, error) {
	var userInfo models.GoogleAuthUser

	configGoogleAPI := &oauth2.Config{
		RedirectURL:  a.config.Auth.GoogleRedirectURL,
		ClientID:     a.config.Auth.GoogleClientID,
		ClientSecret: a.config.Auth.GoogleClientSecret,
		Scopes:       a.config.Auth.GoogleScopes,
		Endpoint:     google.Endpoint,
	}

	token, err := configGoogleAPI.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.config.Auth.GoogleOAuthURL+token.AccessToken, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (a adapter) Auth(ctx context.Context, req *models.AuthRequest) (*models.Session, error) {
	//TODO implement me
	panic("implement me")
}

func (a adapter) Signup(ctx context.Context, req *models.SignupRequest) (*models.Session, error) {
	if req.AuthHash == "" {
		return nil, errors.New("authHash is required")
	}

	_, err := a.db.GetUserByAuthHash(ctx, req.AuthHash)
	if err != nil {
		return nil, err
	}

	res, err := a.user.CreateUser(ctx, &models.CreateUserRequest{
		Username: req.Username,
		Country:  req.Country,
	})
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(res.Id)
	if err != nil {
		return nil, err
	}

	err = a.db.UpdateId(ctx, id, req.AuthHash)
	if err != nil {
		return nil, err
	}

	token, err := a.generateJWT(id)
	if err != nil {
		return nil, err
	}

	err = a.db.RemoveAuthHash(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.Session{
		Token: token,
	}, nil
}
