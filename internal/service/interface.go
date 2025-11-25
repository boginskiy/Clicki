package service

import (
	"net/http"
)

// CrudSrver - interface for standart service.
type CrudSrver interface {
	ReadSetUserURL(*http.Request) ([]byte, error)
	CreateSetURL(*http.Request) ([]byte, error)
	CreateURL(*http.Request) ([]byte, error)
	CheckDB(*http.Request) ([]byte, error)
	ReadURL(*http.Request) ([]byte, error)
	GetHeader() string
}

// CrudSrver - interface for core service.
type CoreSrver interface {
	TakeUserIDFromCtx(*http.Request) int
	EncrypOriginURL() string
}

// CrudSrver - interface for del service.
type DelSrver interface {
	DeleteSetUserURL(req *http.Request) ([]byte, error)
}
