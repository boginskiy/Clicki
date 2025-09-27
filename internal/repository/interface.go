package repository

import (
	"context"
)

type Repository interface {
	Read(context.Context, string) (any, error)
	Create(context.Context, any) (any, error)

	CheckUnic(context.Context, string) bool
	CreateSet(context.Context, any) error
	Ping(context.Context) (bool, error)

	// New
	TakeLastUser(context.Context) int
	ReadSet(context.Context, int) (any, error)
}
