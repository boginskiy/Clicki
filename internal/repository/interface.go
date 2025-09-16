package repository

import (
	"context"
	"database/sql"
)

type URLRepository interface {
	CreateSet(context.Context, any) error
	NewRow(context.Context, string, string, string) any
	Read(context.Context, string) (any, error)
	CheckUnic(context.Context, string) bool
	Create(context.Context, any) error
	GetDB() *sql.DB
}

// Update(any)
// Delete(any)
