package auther

import (
	"errors"
)

var (
	ErrCreateToken    = errors.New(`{"error":"token was not created"}`)
	ErrUserNotFound   = errors.New(`{"error":"user is not found"}`)
	ErrTokenNotValid  = errors.New(`{"error":"token is not valid"}`)
	ErrDataNotValid   = errors.New(`{"error":"data is not valid"}`)
	ErrTokenIsExpired = errors.New(`{"error":"token is expired"}`)
)
