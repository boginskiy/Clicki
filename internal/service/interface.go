package service

import (
	"net/http"
)

type Reader interface {
	ReadSet(*http.Request) ([]byte, error)
	Read(*http.Request) ([]byte, error)
}

type Creater interface {
	CreateSet(*http.Request) ([]byte, error)
	Create(*http.Request) ([]byte, error)
}

type Checker interface {
	CheckDB(*http.Request) ([]byte, error)
}

type Statistician interface {
	GetStats(*http.Request) ([]byte, error)
}

// Servicer is interface for standart service.
type Servicer interface {
	Statistician
	Checker
	Creater
	Reader
}

// DelServicer is interface for del service.
type DelServicer interface {
	DeleteSet(req *http.Request) ([]byte, error)
}
