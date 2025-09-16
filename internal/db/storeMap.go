package db

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	c "github.com/boginskiy/Clicki/cmd/config"
	l "github.com/boginskiy/Clicki/internal/logger"
	m "github.com/boginskiy/Clicki/internal/model"
)

type StoreMap struct {
	Store map[string]string
	mu    sync.RWMutex
}

func NewStoreMap(_ c.VarGetter, _ l.Logger) (*StoreMap, error) {
	return &StoreMap{
		Store: make(map[string]string, SIZE),
	}, nil
}

func (sm *StoreMap) GetDB() *sql.DB {
	return nil
}

func (sm *StoreMap) CloseDB() {
}

func (sm *StoreMap) Read(ctx context.Context, CorrelationID string) (any, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	record, ok := sm.Store[CorrelationID]
	if !ok {
		return nil, errors.New("data is not available")
	}
	return record, nil
}

func (sm *StoreMap) Create(ctx context.Context, record any) error {
	row, ok := record.(*m.URLTb)
	if !ok {
		return errors.New("error in StoreMap>Create")
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Добавление записи в map
	sm.Store[row.CorrelationID] = row.OriginalURL
	return nil
}

func (sm *StoreMap) CheckUnic(ctx context.Context, correlationID string) bool {
	_, ok := sm.Store[correlationID]
	return !ok
}

func (sm *StoreMap) NewRow(ctx context.Context, originURL, shortURL, correlationID string) any {
	return m.NewURLTb(originURL, shortURL, correlationID)
}

func (sf *StoreMap) CreateSet(ctx context.Context, records any) error {
	rows, ok := records.([]m.ResURLSet)
	if !ok || len(rows) == 0 {
		return errors.New("data not valid")
	}

	sf.mu.RLock()

	for _, r := range rows {
		sf.Store[r.CorrelationID] = r.OriginalURL
	}
	sf.mu.RUnlock()
	return nil
}
