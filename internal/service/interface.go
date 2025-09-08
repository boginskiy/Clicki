package service

import (
	"net/http"

	"github.com/boginskiy/Clicki/cmd/config"
)

type CRUDer interface {
	Create(request *http.Request, kwargs config.VarGetter) ([]byte, error)
	CheckPing(request *http.Request) ([]byte, error)
	Read(request *http.Request) ([]byte, error)
}
