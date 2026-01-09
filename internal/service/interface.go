package service

import (
	"context"
	"net/http"

	"github.com/boginskiy/Clicki/internal/protocol"
	p "github.com/boginskiy/Clicki/internal/protocol"
)

type ProtoReader interface {
	Read(ctx context.Context, protocol p.Protocol, request any) ([]byte, error)
	ReadSet(ctx context.Context, protocol p.Protocol) (any, error)
}

type Creater interface {
	Create(*http.Request) ([]byte, error)
}

type ProtoCreater interface {
	Create(ctx context.Context, obj protocol.Protocol, request any) ([]byte, error)
}

type Checker interface {
	CheckDB(*http.Request) ([]byte, error)
}

type Statistician interface {
	GetStats(*http.Request) ([]byte, error)
}

// DelServicer is interface for del service.
type DelServicer interface {
	DeleteSet(req *http.Request) ([]byte, error)
}

// Servicer is interface for standart service.
type Servicer interface {
	Statistician
	ProtoReader
	Checker
	Creater
}

// Servicer is interface for API.
type APIServicer interface {
	ProtoCreater
	Statistician
	ProtoReader
	CreateSet(*http.Request) ([]byte, error)
}
