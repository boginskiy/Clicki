package database

// DataBase - .
type DataBase interface {
	CheckOpen() (bool, error)
	GetDB() any
	CloseDB()
}
