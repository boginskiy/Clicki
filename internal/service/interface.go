package service

import (
	"net/http"
)

// Servicer - interface for standart service.
type Servicer interface {
	ReadSet(*http.Request) ([]byte, error)
	CreateSet(*http.Request) ([]byte, error)
	Read(*http.Request) ([]byte, error)
	Create(*http.Request) ([]byte, error)
	CheckDB(*http.Request) ([]byte, error)
}

// DelServicer - interface for del service.
type DelServicer interface {
	DeleteSet(req *http.Request) ([]byte, error)
}
