package db

import "errors"

const CAP = 10

// Интерфейс Storage описывает операции над хранилищем данных
type Storage interface {
	GetValue(key string) (value string, err error)
	PutValue(key string, value string)
}

type DbStore struct {
	Store map[string]string
}

func NewDbStore() *DbStore {
	return &DbStore{
		Store: make(map[string]string, CAP),
	}
}

func (db *DbStore) GetValue(key string) (value string, err error) {
	value, ok := db.Store[key]
	if !ok {
		return "", errors.New("data is not available")
	}
	return value, nil
}

func (db *DbStore) PutValue(key, value string) {
	db.Store[key] = value
}
