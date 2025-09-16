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
	CheckUnic(string) bool
}
