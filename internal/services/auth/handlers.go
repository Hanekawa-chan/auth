package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"kanji-auth/internal/database"
	"kanji-auth/internal/services/models"
	"net/http"
)

func (a *adapter) GetUserInfoFromGoogleAPI(ctx context.Context, code string) (*models.GoogleAuthUser, error) {
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

func (a *adapter) Auth(ctx context.Context, req *models.AuthRequest) (*models.Session, error) {
	authUser, err := a.getAuthUser(ctx, req)
	if err != nil {
		return nil, err
	}

	var existUser *models.Credentials
	switch req.AuthType.(type) {
	case *models.GoogleAuth:
		user, err := a.db.GetUserByGoogleEmail(ctx, authUser.Email)
		if err != nil {
			return nil, err
		}
		existUser, err = a.db.GetUserByID(ctx, user.ID)
		if err != nil {
			return nil, err
		}
	case *models.PairAuth:
		existUser, err = a.db.GetUserByEmail(ctx, authUser.Email)
		if err != nil {
			return nil, err
		}
	}

	if req.AuthType.(*models.PairAuth) != nil {
		if err == database.ErrNotFound {
			uuID, err := uuid.NewUUID()
			if err != nil {
				return nil, err
			}
			authHash, err := generateAuthHash()
			if err != nil {
				return nil, err
			}
			hash, err := bcrypt.GenerateFromPassword([]byte(authUser.Password), 12)
			if err != nil {
				return nil, err
			}

			user := &models.Credentials{
				ID:       uuID,
				Email:    authUser.Email,
				Password: string(hash),
				AuthHash: authHash,
			}

			if err := a.db.CreateUser(ctx, user); err != nil {
				return nil, err
			}

			return &models.Session{
				AuthHash: authHash,
			}, nil
		}
		if err != nil {
			return nil, err
		}
	}

	if len(existUser.AuthHash) > 0 {
		return &models.Session{
			AuthHash: existUser.AuthHash,
		}, nil
	}

	token, err := a.generateJWT(existUser.ID)
	if err != nil {
		return nil, err
	}

	return &models.Session{
		Token: token,
	}, nil
}

func (a *adapter) Signup(ctx context.Context, req *models.SignupRequest) (*models.Session, error) {
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

func (a *adapter) getAuthUser(ctx context.Context, req *models.AuthRequest) (*models.Credentials, error) {
	switch req.AuthType.(type) {
	case *models.GoogleAuth:
		v := req.AuthType.(*models.GoogleAuth)
		return a.getUserByGoogle(ctx, v)
	case *models.PairAuth:
		v := req.AuthType.(*models.PairAuth)
		return a.getUserByPair(ctx, v.Email, v.Password)
	}
	return nil, errors.New("invalid auth_type")
}

func (a *adapter) getUserByGoogle(ctx context.Context, req *models.GoogleAuth) (*models.Credentials, error) {
	googleUser, err := a.GetUserInfoFromGoogleAPI(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &models.Credentials{Email: googleUser.Email}, nil
}

func (a *adapter) getUserByPair(ctx context.Context, login string, password string) (*models.Credentials, error) {
	err := a.validateEmail(login)
	if err != nil {
		return nil, errors.New("email isn't valid")
	}

	err = a.validatePassword(password)
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	user, err := a.db.GetUserByEmail(ctx, login)
	if err == database.ErrNotFound {
		return &models.Credentials{
			Email:    login,
			Password: string(hash),
		}, nil
	}

	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), hash); err != nil {
		return nil, errors.New("password is wrong")
	}
	return user, nil
}
