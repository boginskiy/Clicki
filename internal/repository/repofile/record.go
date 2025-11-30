package repofile

import (
	"context"
	"encoding/json"

	"github.com/boginskiy/Clicki/internal/errs"
	"github.com/boginskiy/Clicki/internal/model"
)

type FileRecordRepo struct {
	Repo *RepoFile
}

func NewFileRecordRepo(repo *RepoFile) *FileRecordRepo {
	return &FileRecordRepo{
		Repo: repo,
	}
}

// CheckUniqueRecord.
func (r *FileRecordRepo) CheckUniqueRecord(ctx context.Context, correlID string) bool {
	_, ok := r.Repo.tmpStore[correlID]
	return !ok
}

// ReadLastRecord.
func (r *FileRecordRepo) ReadLastRecord(ctx context.Context) int {
	return r.Repo.lastUser
}

// ReadRecord.
func (r *FileRecordRepo) ReadRecord(ctx context.Context, correlID string) (any, error) {
	r.Repo.muR.RLock()
	defer r.Repo.muR.RUnlock()

	record, ok := r.Repo.tmpStore[correlID]
	if !ok {
		return nil, errs.NewErrPlace("data is not available", nil)
	}
	return record, nil
}

// CreateRecord.
func (r *FileRecordRepo) CreateRecord(ctx context.Context, preRecord any) (any, error) {
	row, ok := preRecord.(*model.URLTb)
	if !ok {
		return nil, errs.NewErrPlace("type is not available", nil)
	}

	// Логика, если данные уже есть в Store.
	r.Repo.muR.RLock()
	if correlID, ok := r.Repo.uniqueFields[row.OriginalURL]; ok {
		return r.Repo.tmpStore[correlID], errs.ErrUniqueData
	}
	r.Repo.muR.RUnlock()

	// Логика, если данные отсутствуют в Store.
	r.Repo.mu.Lock()
	r.Repo.lastRecord += 1
	row.ID = r.Repo.lastRecord
	r.Repo.tmpStore[row.CorrelationID] = row
	r.Repo.uniqueFields[row.OriginalURL] = row.CorrelationID
	r.Repo.mu.Unlock()

	jsonData, err := json.Marshal(row)
	if err != nil {
		return nil, errs.NewErrPlace("type is not available", err)
	}
	jsonData = append(jsonData, byte('\n'))

	_, err = r.Repo.Store.Write(jsonData)
	return row, err
}
