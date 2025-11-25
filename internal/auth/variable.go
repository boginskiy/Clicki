package auth

import (
	"errors"
)

// Errors for this package
var (
	ErrCreateToken    = errors.New(`{"error":"token was not created"}`)
	ErrUserNotFound   = errors.New(`{"error":"user is not found"}`)
	ErrTokenNotValid  = errors.New(`{"error":"token is not valid"}`)
	ErrDataNotValid   = errors.New(`{"error":"data is not valid"}`)
	ErrTokenIsExpired = errors.New(`{"error":"token is expired"}`)
)

// contextKey - is type of key for values for context request.
type contextKey struct{}

// CtxUserID - is key for userID to save.
var CtxUserID = contextKey{}
