package db

type DBer interface {
	CheckOpen() (bool, error)
	GetDB() any
	CloseDB()
}
