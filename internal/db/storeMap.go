package db

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	c "github.com/boginskiy/Clicki/cmd/config"
	l "github.com/boginskiy/Clicki/internal/logger"
	m "github.com/boginskiy/Clicki/internal/model"
)

type StoreMap struct {
	Store map[string]*m.URLTb
	mu    sync.RWMutex
}

func NewStoreMap(_ c.VarGetter, _ l.Logger) (*StoreMap, error) {
	return &StoreMap{
		Store: make(map[string]*m.URLTb, SIZE),
	}, nil
}

func (sm *StoreMap) GetDB() *sql.DB {
	return nil
}

func (sm *StoreMap) CloseDB() {
}

func (sm *StoreMap) Read(shortURL string) (any, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	record, ok := sm.Store[shortURL]
	if !ok {
		fmt.Println(">>", shortURL)
		fmt.Println(">>", sm.Store)
		return nil, errors.New("data is not available")
	}
	return record, nil
}

func (sm *StoreMap) Create(record any) error {
	row, ok := record.(*m.URLTb)
	if !ok {
		return errors.New("error in StoreMap>Create")
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Добавление записи в map
	sm.Store[row.ShortURL] = row
	return nil
}

func (sm *StoreMap) CheckUnic(shortURL string) bool {
	_, ok := sm.Store[shortURL]
	return !ok
}

func (sm *StoreMap) NewRow(originURL, shortURL string) any {
	return m.NewURLTb(originURL, shortURL)
}
