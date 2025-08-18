package db

import (
	"errors"
	"sync"
)

const CAP = 10

// Интерфейс Storage описывает операции над хранилищем данных
type Storage interface {
	GetValue(key string) (value string, err error)
	PutValue(key string, value string)
}

type DBStore struct {
	mu    sync.RWMutex
	Store map[string]string
}

func NewDBStore() *DBStore {
	return &DBStore{
		Store: make(map[string]string, CAP),
	}
}

func (db *DBStore) GetValue(key string) (value string, err error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	value, ok := db.Store[key]
	if !ok {
		return "", errors.New("data is not available")
	}
	return value, nil
}

func (db *DBStore) PutValue(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.Store[key] = value
}
