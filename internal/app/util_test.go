package app

import (
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"os"
	"testing"
	"time"
)

func Test_service_generateAccessToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		logger       *zerolog.Logger
		jwtGenerator JWT
		db           Database
		user         User
		api          Api
	}
	type args struct {
		userID uuid.UUID
	}

	logger := zerolog.New(os.Stdout)
	zl := &logger
	f := fields{
		logger:       zl,
		jwtGenerator: NewMockJWT(ctrl),
		db:           NewMockDatabase(ctrl),
		user:         NewMockUser(ctrl),
		api:          NewMockApi(ctrl),
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func(jwt *MockJWT, a args)
		want       string
		wantErr    bool
	}{
		{
			name:   "success generate access token",
			fields: f,
			args:   args{userID: uuid.MustParse("19e64cf1-7a02-4504-9690-cb81d35b8375")},
			beforeTest: func(jwt *MockJWT, a args) {
				claims := make(map[string]interface{})
				claims["user_id"] = a.userID
				claims["exp"] = time.Now().Add(AccessExp).Unix()
				jwt.EXPECT().Generate(claims).Return("test", nil)
			},
			want:    "test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &service{
				logger:       tt.fields.logger,
				jwtGenerator: tt.fields.jwtGenerator,
				db:           tt.fields.db,
				user:         tt.fields.user,
				api:          tt.fields.api,
			}
			if tt.beforeTest != nil {
				tt.beforeTest(a.jwtGenerator.(*MockJWT), tt.args)
			}
			got, err := a.generateAccessToken(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateAccessToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_generateRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		logger       *zerolog.Logger
		jwtGenerator JWT
		db           Database
		user         User
		api          Api
	}

	logger := zerolog.New(os.Stdout)
	zl := &logger
	f := fields{
		logger:       zl,
		jwtGenerator: NewMockJWT(ctrl),
		db:           NewMockDatabase(ctrl),
		user:         NewMockUser(ctrl),
		api:          NewMockApi(ctrl),
	}

	tests := []struct {
		name       string
		fields     fields
		beforeTest func(jwt *MockJWT)
		want       string
		wantErr    bool
	}{
		{
			name:   "success generate refresh token",
			fields: f,
			beforeTest: func(jwt *MockJWT) {
				claims := make(map[string]interface{})
				claims["exp"] = time.Now().Add(RefreshExp).Unix()
				jwt.EXPECT().Generate(claims).Return("test", nil)
			},
			want:    "test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &service{
				logger:       tt.fields.logger,
				jwtGenerator: tt.fields.jwtGenerator,
				db:           tt.fields.db,
				user:         tt.fields.user,
				api:          tt.fields.api,
			}
			if tt.beforeTest != nil {
				tt.beforeTest(a.jwtGenerator.(*MockJWT))
			}
			got, err := a.generateRefreshToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateRefreshToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_parseRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		logger       *zerolog.Logger
		jwtGenerator JWT
		db           Database
		user         User
		api          Api
	}
	type args struct {
		token string
	}
	logger := zerolog.New(os.Stdout)
	zl := &logger
	f := fields{
		logger:       zl,
		jwtGenerator: NewMockJWT(ctrl),
		db:           NewMockDatabase(ctrl),
		user:         NewMockUser(ctrl),
		api:          NewMockApi(ctrl),
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		beforeTest func(jwt *MockJWT, a args)
		wantErr    bool
	}{
		{
			name:   "success parse token",
			fields: f,
			args:   args{token: "token"},
			beforeTest: func(jwt *MockJWT, a args) {
				jwt.EXPECT().ParseToken(a.token).Return(nil, nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &service{
				logger:       tt.fields.logger,
				jwtGenerator: tt.fields.jwtGenerator,
				db:           tt.fields.db,
				user:         tt.fields.user,
				api:          tt.fields.api,
			}
			if tt.beforeTest != nil {
				tt.beforeTest(a.jwtGenerator.(*MockJWT), tt.args)
			}
			if err := a.parseRefreshToken(tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("parseRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_validateEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		logger       *zerolog.Logger
		jwtGenerator JWT
		db           Database
		user         User
		api          Api
	}
	type args struct {
		email string
	}
	logger := zerolog.New(os.Stdout)
	zl := &logger
	f := fields{
		logger:       zl,
		jwtGenerator: NewMockJWT(ctrl),
		db:           NewMockDatabase(ctrl),
		user:         NewMockUser(ctrl),
		api:          NewMockApi(ctrl),
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "success validate email",
			fields:  f,
			args:    args{email: "test@mail.ru"},
			wantErr: false,
		},
		{
			name:    "failed validate email without left side",
			fields:  f,
			args:    args{email: "@test.ru"},
			wantErr: true,
		},
		{
			name:    "failed validate email without mail",
			fields:  f,
			args:    args{email: "test@.ru"},
			wantErr: true,
		},
		{
			name:    "failed validate email without right side",
			fields:  f,
			args:    args{email: "test@"},
			wantErr: true,
		},
		{
			name:    "failed validate email without @ and right side",
			fields:  f,
			args:    args{email: "test"},
			wantErr: true,
		},
		{
			name:    "failed validate email without everything",
			fields:  f,
			args:    args{email: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &service{
				logger:       tt.fields.logger,
				jwtGenerator: tt.fields.jwtGenerator,
				db:           tt.fields.db,
				user:         tt.fields.user,
				api:          tt.fields.api,
			}
			if err := a.validateEmail(tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("validateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_validatePair(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		logger       *zerolog.Logger
		jwtGenerator JWT
		db           Database
		user         User
		api          Api
	}
	type args struct {
		login    string
		password string
	}
	logger := zerolog.New(os.Stdout)
	zl := &logger
	f := fields{
		logger:       zl,
		jwtGenerator: NewMockJWT(ctrl),
		db:           NewMockDatabase(ctrl),
		user:         NewMockUser(ctrl),
		api:          NewMockApi(ctrl),
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "success validate email",
			fields:  f,
			args:    args{login: "test@test.ru", password: "password"},
			wantErr: false,
		},
		{
			name:    "invalid validate password",
			fields:  f,
			args:    args{login: "test@test.ru", password: "passwor"},
			wantErr: true,
		},
		{
			name:    "invalid validate login",
			fields:  f,
			args:    args{login: "test@test.", password: "password"},
			wantErr: true,
		},
		{
			name:    "invalid validate both",
			fields:  f,
			args:    args{login: "test@test.", password: "passwor"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &service{
				logger:       tt.fields.logger,
				jwtGenerator: tt.fields.jwtGenerator,
				db:           tt.fields.db,
				user:         tt.fields.user,
				api:          tt.fields.api,
			}
			if err := a.validatePair(tt.args.login, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("validatePair() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_validatePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		logger       *zerolog.Logger
		jwtGenerator JWT
		db           Database
		user         User
		api          Api
	}
	type args struct {
		password string
	}

	logger := zerolog.New(os.Stdout)
	zl := &logger
	f := fields{
		logger:       zl,
		jwtGenerator: NewMockJWT(ctrl),
		db:           NewMockDatabase(ctrl),
		user:         NewMockUser(ctrl),
		api:          NewMockApi(ctrl),
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "success validate email",
			fields:  f,
			args:    args{password: "password"},
			wantErr: false,
		},
		{
			name:    "invalid validate email",
			fields:  f,
			args:    args{password: "passwor"},
			wantErr: true,
		},
		{
			name:    "invalid validate email",
			fields:  f,
			args:    args{password: "passworklnagkn;dagn;;nagonaeoneagonagonagronagnagnkegrkaglkanklklkfsgkmfmkegk"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &service{
				logger:       tt.fields.logger,
				jwtGenerator: tt.fields.jwtGenerator,
				db:           tt.fields.db,
				user:         tt.fields.user,
				api:          tt.fields.api,
			}
			if err := a.validatePassword(tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("validatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
