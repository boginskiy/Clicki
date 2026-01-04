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
func (r *FileRecordRepo) CreateRecord(ctx context.Context, record *model.URLTb) (any, error) {
	// Logic, if data in a Store.
	r.Repo.muR.RLock()
	if correlID, ok := r.Repo.uniqueFields[record.OriginalURL]; ok {
		return r.Repo.tmpStore[correlID], errs.ErrUniqueData
	}
	r.Repo.muR.RUnlock()

	// Logic, if data not in a Store.
	r.Repo.mu.Lock()
	r.Repo.lastRecord += 1
	record.ID = r.Repo.lastRecord
	r.Repo.tmpStore[record.CorrelationID] = record
	r.Repo.uniqueFields[record.OriginalURL] = record.CorrelationID
	r.Repo.mu.Unlock()

	jsonData, err := json.Marshal(record)
	if err != nil {
		return nil, errs.NewErrPlace("type is not available", err)
	}
	jsonData = append(jsonData, byte('\n'))

	_, err = r.Repo.Store.Write(jsonData)
	return record, err
}
