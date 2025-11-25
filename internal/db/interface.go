package db

// DataBase - .
type DataBase interface {
	CheckOpen() (bool, error)
	GetDB() any
	CloseDB()
}
