package db2

import "database/sql"

type DBConnecter interface {
	GetDB() *sql.DB
	CloseDB()
}
