package app

import (
	"context"
	"encoding/json"
	"github.com/Hanekawa-chan/kanji-auth/internal/services/errors"
	"github.com/Hanekawa-chan/kanji-auth/internal/services/models"
	kanjiJwt "github.com/Hanekawa-chan/kanji-jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"net/http"
)

func (a *service) getUserInfoFromGoogleAPI(ctx context.Context, code string) (*models.GoogleAuthUser, error) {
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

func (a *service) Auth(ctx context.Context, req *models.AuthRequest) (*models.Session, error) {
	authUser, err := a.getAuthUser(ctx, req)
	if err != nil {
		return nil, err
	}

	var existUser *models.Credentials
	switch req.AuthType.(type) {
	case *models.GoogleAuth:
		user, err := a.db.GetUserByGoogleEmail(ctx, authUser.Login)
		if err != nil {
			return nil, err
		}
		existUser, err = a.db.GetUserByID(ctx, user.Id)
		if err != nil {
			return nil, err
		}
	case *models.PairAuth:
		existUser, err = a.db.GetUserByEmail(ctx, authUser.Login)
		if err != nil {
			return nil, err
		}
	}

	if req.AuthType.(*models.PairAuth) != nil {
		if err == errors.ErrNotFound {
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
				Id:       uuID,
				Login:    authUser.Login,
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

	token, err := a.generateJWT(existUser.Id)
	if err != nil {
		return nil, err
	}

	return &models.Session{
		Token: token,
	}, nil
}

func (a *service) Signup(ctx context.Context, req *models.SignupRequest) (*models.Session, error) {
	if req.AuthHash == "" {
		return nil, errors.ErrEmptyRequired
	}

	_, err := a.db.GetUserByAuthHash(ctx, req.AuthHash)
	if err != nil {
		return nil, err
	}

	res, err := a.user.CreateUser(ctx, &models.CreateUserRequest{
		Email:   req.Email,
		Country: req.Country,
	})
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(res)
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

func (a *service) Link(ctx context.Context, req *models.AuthRequest) error {
	id, err := kanjiJwt.GetUserId(ctx, a.jwtGenerator.(*kanjiJwt.Generator))
	if err != nil {
		return err
	}
	switch req.AuthType.(type) {
	case *models.GoogleAuth:
		v := req.AuthType.(*models.GoogleAuth)
		cred, err := a.linkGoogle(ctx, v, id)
		if err != nil {
			return err
		}
		err = a.db.CreateGoogle(ctx, cred)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.ErrType
}

func (a *service) getAuthUser(ctx context.Context, req *models.AuthRequest) (*models.Credentials, error) {
	switch req.AuthType.(type) {
	case *models.GoogleAuth:
		v := req.AuthType.(*models.GoogleAuth)
		return a.getUserByGoogle(ctx, v)
	case *models.PairAuth:
		v := req.AuthType.(*models.PairAuth)
		return a.getUserByPair(ctx, v.Email, v.Password)
	}
	return nil, errors.ErrType
}

func (a *service) getUserByGoogle(ctx context.Context, req *models.GoogleAuth) (*models.Credentials, error) {
	googleUser, err := a.getUserInfoFromGoogleAPI(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &models.Credentials{Login: googleUser.Email}, nil
}

func (a *service) getUserByPair(ctx context.Context, login string, password string) (*models.Credentials, error) {
	err := a.validateEmail(login)
	if err != nil {
		return nil, errors.ErrValidation
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
	if err == errors.ErrNotFound {
		return &models.Credentials{
			Login:    login,
			Password: string(hash),
		}, nil
	}

	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), hash); err != nil {
		return nil, errors.ErrValidation
	}
	return user, nil
}

func (a *service) linkGoogle(ctx context.Context, req *models.GoogleAuth, id uuid.UUID) (*models.Google, error) {
	googleUser, err := a.getUserInfoFromGoogleAPI(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &models.Google{Id: id, Email: googleUser.Email, GoogleId: googleUser.ID}, nil
}
