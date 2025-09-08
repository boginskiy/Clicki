package db

import (
	"errors"
	"sync"
)

const CAP = 10

type DBStore struct {
	mu         sync.RWMutex
	Store      map[string]string
	fileWorker FileWorker
}

func NewDBStore(fileWorker FileWorker) *DBStore {
	return &DBStore{
		Store:      fileWorker.DataRecovery(&StoreModel{}),
		fileWorker: fileWorker,
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

	// Добавление записи в map
	db.Store[key] = value

	// Добавление записи в файл
	nextLine := db.fileWorker.GetNextLine()
	db.fileWorker.Save(NewStoreModel(nextLine, key, value))
}
