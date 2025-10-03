package service

import (
	"net/http"
)

type CrudSrver interface {
	ReadSetUserURL(*http.Request) ([]byte, error)
	CreateSetURL(*http.Request) ([]byte, error)
	CreateURL(*http.Request) ([]byte, error)
	CheckDB(*http.Request) ([]byte, error)
	ReadURL(*http.Request) ([]byte, error)
	GetHeader() string
}

type CoreSrver interface {
	TakeUserIDFromCtx(*http.Request) int
	EncrypOriginURL() string
}

type DelSrver interface {
	DeleteSetUserURL(req *http.Request) ([]byte, error)
}
