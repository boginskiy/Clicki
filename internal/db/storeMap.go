package db

import (
	"context"
	"database/sql"
	"sync"

	c "github.com/boginskiy/Clicki/cmd/config"
	e "github.com/boginskiy/Clicki/internal/errors"
	l "github.com/boginskiy/Clicki/internal/logger"
	m "github.com/boginskiy/Clicki/internal/model"
)

type StoreMap struct {
	Store        map[string]*m.URLTb
	uniqueFields map[string]string
	mu           sync.Mutex
	muR          sync.RWMutex
}

func NewStoreMap(_ c.VarGetter, _ l.Logger) (*StoreMap, error) {
	return &StoreMap{
		Store:        make(map[string]*m.URLTb, SIZE),
		uniqueFields: make(map[string]string, SIZE),
	}, nil
}

func (sm *StoreMap) GetDB() *sql.DB {
	return nil
}

func (sm *StoreMap) CloseDB() {
}

func (sm *StoreMap) CheckUnic(ctx context.Context, correlID string) bool {
	_, ok := sm.Store[correlID]
	return !ok
}

func (sm *StoreMap) Read(ctx context.Context, correlID string) (any, error) {
	sm.muR.RLock()
	defer sm.muR.RUnlock()

	record, ok := sm.Store[correlID]
	if !ok {
		return nil, e.NewErrPlace("data is not available", nil)
	}
	return record, nil
}

func (sm *StoreMap) Create(ctx context.Context, preRecord any) (any, error) {
	row, ok := preRecord.(*m.URLTb)
	if !ok {
		return nil, e.NewErrPlace("data not valid", nil)
	}

	// Логика, если данные уже есть в Store
	sm.muR.RLock()
	defer sm.muR.RUnlock()

	if correlID, ok := sm.uniqueFields[row.OriginalURL]; ok {
		return sm.Store[correlID], e.ErrUniqueData
	}

	// Добавление записи в map
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.Store[row.CorrelationID] = row
	sm.uniqueFields[row.OriginalURL] = row.CorrelationID

	return row, nil
}

func (sm *StoreMap) CreateSet(ctx context.Context, records any) error {
	rows, ok := records.([]m.ResURLSet)
	if !ok || len(rows) == 0 {
		return e.NewErrPlace("data not valid", nil)
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	for _, r := range rows {

		// TODO! Перекладка с ResURLSet в URLTb не супер оптимально
		// однако пока так ...

		sm.Store[r.CorrelationID] = &m.URLTb{
			ID:            0,
			OriginalURL:   r.OriginalURL,
			ShortURL:      r.ShortURL,
			CorrelationID: r.CorrelationID,
			CreatedAt:     r.CreatedAt,
		}

		sm.uniqueFields[r.OriginalURL] = r.CorrelationID
	}
	return nil
}
