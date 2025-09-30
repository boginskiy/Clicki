package service

import (
	"net/http"
)

type CRUDer interface {
	DeleteSetUserURL(*http.Request) ([]byte, error)
	ReadSetUserURL(*http.Request) ([]byte, error)
	CreateSetURL(*http.Request) ([]byte, error)
	CreateURL(*http.Request) ([]byte, error)
	CheckDB(*http.Request) ([]byte, error)
	ReadURL(*http.Request) ([]byte, error)
	GetHeader() string
}

type CoreServicer interface {
	takeUserIDFromCtx(*http.Request) int
	encrypOriginURL() string
}
