package repository

import (
	"context"
	"database/sql"
)

type URLRepository interface {
	Read(context.Context, string) (any, error)
	CheckUnic(context.Context, string) bool
	CreateSet(context.Context, any) error
	Create(context.Context, any) error
	GetDB() *sql.DB
}

// Update(any)
// Delete(any)
