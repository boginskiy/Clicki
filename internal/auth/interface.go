package auth

import (
	"context"
	"net/http"
)

// Auther - for authentication of users.
type Auther interface {
	Authentication(req *http.Request) (*http.Cookie, int, error)
}

type AutherGRPC interface {
	Authentication(ctx context.Context) (string, int, error)
}

// JWTer - for a work of JWT authentication.
type JWTer interface {
	GetIDAndValidJWT(tokenStr string) (int, error)
	CreateJWT(userID int) (string, error)
}
