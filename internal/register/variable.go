package register

import (
	"errors"
	"time"
)

var (
	TokenNotValid  = errors.New(`{"error":"token is not valid"}`)
	DataNotValid   = errors.New(`{"error":"data not valid"}`)
	TokenIsExpired = errors.New(`{"error":"token is expired"}`)
)

var (
	MessAuthCompleted   = []byte(`{"mess":"вы прошли аутентификацию, повторите запрос"}`)
	MessAuthCompleted2  = []byte(`{"mess":"вы прошли повторную аутентификацию, повторите запрос"}`)
	MessRegistCompleted = []byte(`{"mess":"вы прошли регистрацию, повторите запрос"}`)
	MessRegistIsBad     = []byte(`{"mess":"вы не прошли регистрацию. Повторите запрос позже"}`)
	MessProcessingIsBad = []byte(`{"mess":"что-то пошло не так... Повторите запрос позже"}`)
	MessAuthIsBad       = []byte(`{"mess":"вы не прошли аутентификацию, повторите запрос"}`)
	MessNeedRegist      = []byte(`{"mess":"вам необходимо зарегистрироваться для просмотра ресурса"}`)
	MessGoodUser        = []byte(`{"mess":"you are good boys!"}`)
)

var EmptyByteSlice = []byte{}

// TODO! Вынести в переменные окружения
const TOKEN_EXP = time.Second * 10
const SECRET_KEY = "1234"
const NAME_COKI = "auth_token"
