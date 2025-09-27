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
const TOKENEXP = time.Second * 10
const SECRETKEY = "1234"
