package repository

import (
	"context"
	"sync"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/db"
	cerr "github.com/boginskiy/Clicki/internal/error"
	mod "github.com/boginskiy/Clicki/internal/model"
)

type RepositoryMapURL struct {
	Kwargs conf.VarGetter
	DB     db.DBer // map[string]*mod.URLTb

	store        map[string]*mod.URLTb
	uniqueFields map[string]string
	muR          sync.RWMutex
	mu           sync.Mutex
}

func NewRepositoryMapURL(kwargs conf.VarGetter, dber db.DBer) (Repository, error) {
	store, ok := dber.GetDB().(map[string]*mod.URLTb)
	if !ok {
		return nil, cerr.NewErrPlace("database not valid", nil)
	}

	return &RepositoryMapURL{
		Kwargs:       kwargs,
		DB:           dber,
		uniqueFields: make(map[string]string, SIZE),
		store:        store,
	}, nil
}

// ReadLastRecord - Для реализации interface
func (rm *RepositoryMapURL) ReadLastRecord(ctx context.Context) int {
	return 0
}

// MarkerRecords - Для реализации interface
func (rm *RepositoryMapURL) MarkerRecords(ctx context.Context, messages ...DelMessage) error {
	return nil
}

// DeleteRecords - Для реализации interface
func (rm *RepositoryMapURL) DeleteRecords(ctx context.Context) error {
	return nil
}

// PingDB - Для реализации interface
func (rm *RepositoryMapURL) PingDB(ctx context.Context) (bool, error) {
	return rm.DB.CheckOpen()
}

func (rm *RepositoryMapURL) CheckUnicRecord(ctx context.Context, correlID string) bool {
	_, ok := rm.store[correlID]
	return !ok
}

func (rm *RepositoryMapURL) ReadRecord(ctx context.Context, correlID string) (any, error) {
	rm.muR.RLock()
	defer rm.muR.RUnlock()

	record, ok := rm.store[correlID]
	if !ok {
		return nil, cerr.NewErrPlace("data is not available", nil)
	}
	return record, nil
}

func (rm *RepositoryMapURL) CreateRecord(ctx context.Context, preRecord any) (any, error) {
	row, ok := preRecord.(*mod.URLTb)
	if !ok {
		return nil, cerr.NewErrPlace("data not valid", nil)
	}

	// Логика, если данные уже есть в Store
	rm.muR.RLock()
	defer rm.muR.RUnlock()

	if correlID, ok := rm.uniqueFields[row.OriginalURL]; ok {
		return rm.store[correlID], cerr.ErrUniqueData
	}

	// Добавление записи в map
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.store[row.CorrelationID] = row
	rm.uniqueFields[row.OriginalURL] = row.CorrelationID

	return row, nil
}

func (rm *RepositoryMapURL) CreateRecords(ctx context.Context, records any) error {
	rows, ok := records.([]mod.ResURLSet)
	if !ok || len(rows) == 0 {
		return cerr.NewErrPlace("data not valid", nil)
	}

	rm.mu.Lock()
	defer rm.mu.Unlock()

	for _, r := range rows {

		// TODO! Перекладка с ResURLSet в URLTb не супер оптимально
		// однако пока так ...

		rm.store[r.CorrelationID] = &mod.URLTb{
			ID:            0,
			OriginalURL:   r.OriginalURL,
			ShortURL:      r.ShortURL,
			CorrelationID: r.CorrelationID,
			CreatedAt:     r.CreatedAt,
			UserID:        r.UserID,
		}

		rm.uniqueFields[r.OriginalURL] = r.CorrelationID
	}
	return nil
}

func (rm *RepositoryMapURL) ReadRecords(ctx context.Context, userID int) (any, error) {
	records := []mod.ResUserURLSet{}

	for _, v := range rm.store {
		if v.UserID == userID {
			records = append(records, mod.ResUserURLSet{
				OriginalURL: v.OriginalURL,
				ShortURL:    v.ShortURL})
		}
	}
	return records, nil
}
