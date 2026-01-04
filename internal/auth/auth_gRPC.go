package auth

import (
	"context"
	"errors"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
	repo "github.com/boginskiy/Clicki/internal/repository"
	"google.golang.org/grpc/metadata"
)

type AuthGRPC struct {
	Cfg        conf.Config
	Logg       logg.Logger
	Repo       repo.Repository
	LastUser   int
	JWTService JWTer
}

func NewAuthGRPC(config conf.Config, logger logg.Logger, repo repo.Repository) *AuthGRPC {
	return &AuthGRPC{
		Cfg:        config,
		Logg:       logger,
		Repo:       repo,
		JWTService: NewJWTService(config),

		// Louding last userID from last record.
		LastUser: repo.ReadLastRecord(context.TODO()),
	}
}

func (a *AuthGRPC) Authorization(ctx context.Context) (string, int, error) {
	UserID := a.nextUser()
	token, err := a.JWTService.CreateJWT(UserID)
	if err != nil {
		a.Logg.RaiseError(err, "AuthGRPC>Authorization>CreateJWT", nil)
		return "", 0, ErrCreateToken
	}
	return token, UserID, nil
}

func (a *AuthGRPC) Authentication(ctx context.Context) (string, int, error) {
	var token string

	// Take token from context.
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		val := md.Get("authorization")
		if len(val) > 0 {
			token = val[0]
		}
	}

	// Authorization. There is not token.
	if len(token) == 0 {
		return a.Authorization(ctx)
	}

	// Authentication. There is token.
	UserID, err := a.JWTService.GetIDAndValidJWT(token)

	// Check token. If User was not found.
	if UserID <= 0 {
		a.Logg.RaiseInfo(ErrUserNotFound.Error(), logg.Fields{"userID": UserID})
		return "", 0, ErrUserNotFound
	}

	// Check token. If token is bad.
	if err != nil {
		return a.RefreshToken(UserID, err)
	}
	return token, UserID, nil
}

func (a *AuthGRPC) RefreshToken(userID int, err error) (string, int, error) {
	if errors.Is(err, ErrTokenIsExpired) || errors.Is(err, ErrTokenNotValid) {
		token, err := a.JWTService.CreateJWT(userID)
		if err != nil {
			a.Logg.RaiseError(err, "problem with creating token", nil)
			return "", 0, ErrCreateToken
		}
		// Give fresh token.
		return token, userID, nil
	}
	a.Logg.RaiseError(err, "update token is impossible", nil)
	return "", 0, ErrTokenNotValid
}

func (a *AuthGRPC) nextUser() int {
	a.LastUser += 1
	return a.LastUser
}
