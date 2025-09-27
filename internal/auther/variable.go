package auther

import (
	"errors"
	"time"
)

var (
	ErrTokenNotValid  = errors.New(`{"error":"token is not valid"}`)
	ErrDataNotValid   = errors.New(`{"error":"data not valid"}`)
	ErrTokenIsExpired = errors.New(`{"error":"token is expired"}`)
)

// TODO! Вынести в переменные окружения
const TOKEN_EXP = time.Second * 300
const SECRET_KEY = "1234"
