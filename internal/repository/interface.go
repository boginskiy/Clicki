package repository

import (
	"context"
)

type Repository interface {
	// URL
	Read(context.Context, string) (any, error)
	Create(context.Context, any) (any, error)
	CheckUnic(context.Context, string) bool
	CreateSet(context.Context, any) error
	Ping(context.Context) (bool, error)
	ReadSet(context.Context, int) (any, error)

	// User
	CreateUser(context.Context, any) (int, error)
	ReadUser(context.Context, int) any
}

type RepositoryURL interface {
	Read(context.Context, string) (any, error)
	Create(context.Context, any) (any, error)
	CheckUnic(context.Context, string) bool
	CreateSet(context.Context, any) error
	Ping(context.Context) (bool, error)
	ReadSet(context.Context, int) (any, error)
}

type RepositoryUser interface {
	CreateUser(context.Context, any) (int, error)
	ReadUser(context.Context, int) any
}
