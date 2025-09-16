package service

import (
	"net/http"
)

type CRUDer interface {
	CheckPing(request *http.Request) ([]byte, error)
	SetBatch(request *http.Request) ([]byte, error)
	Create(request *http.Request) ([]byte, error)
	Read(request *http.Request) ([]byte, error)
	GetHeader() string
}
