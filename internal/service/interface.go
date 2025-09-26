package service

import (
	"net/http"
)

type CRUDer interface {
	CheckDB(request *http.Request) ([]byte, error)
	CreateSetURL(request *http.Request) ([]byte, error)
	ReadSetUserURL(request *http.Request) ([]byte, error)
	CreateURL(request *http.Request) ([]byte, error)
	ReadURL(request *http.Request) ([]byte, error)
	GetHeader() string
}

type CoreServicer interface {
	takeUserIDFromCtx(*http.Request) int
	encrypOriginURL() string
}
