package repository

import (
	"database/sql"
)

type URLRepository interface {
	Create(any) error
	Read(string) (any, error)
	// Update(any)
	// Delete(any)
	NewRow(string, string) any
	CheckUnic(string) bool
	GetDB() *sql.DB
}
