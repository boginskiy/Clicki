package db

import (
	"database/sql"
)

type DBer interface {
	GetDB() *sql.DB
	CloseDB()
}

type Tbler interface {
	Create(any) error
	Read(string) (any, error)
	// NewRow(string, string) any
	CheckUnic(string) bool
}
