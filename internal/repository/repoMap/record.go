package repoMap

import (
	"context"

	"github.com/boginskiy/Clicki/internal/errs"
	"github.com/boginskiy/Clicki/internal/model"
)

type MapRecordRepo struct {
	Repo *RepoMap
}

func NewMapRecordRepo(repo *RepoMap) *MapRecordRepo {
	return &MapRecordRepo{
		Repo: repo,
	}
}

// ReadLastRecord - for interface.
func (r *MapRecordRepo) ReadLastRecord(ctx context.Context) int {
	return 0
}

// CheckUniqueRecord - for interface.
func (r *MapRecordRepo) CheckUniqueRecord(ctx context.Context, correlID string) bool {
	_, ok := r.Repo.Store[correlID]
	return !ok
}

// ReadRecord - for interface.
func (r *MapRecordRepo) ReadRecord(ctx context.Context, correlID string) (any, error) {
	r.Repo.muR.RLock()
	defer r.Repo.muR.RUnlock()

	record, ok := r.Repo.Store[correlID]
	if !ok {
		return nil, errs.NewErrPlace("data is not available in ReadRecord", nil)
	}
	return record, nil
}

// CreateRecord - for interface.
func (r *MapRecordRepo) CreateRecord(ctx context.Context, records any) (any, error) {
	row, ok := records.(*model.URLTb)
	if !ok {
		return nil, errs.NewErrPlace("records in CreateRecord not valid", nil)
	}

	// If data is in the Store.
	r.Repo.muR.RLock()
	defer r.Repo.muR.RUnlock()

	if correlID, ok := r.Repo.uniqueFields[row.OriginalURL]; ok {
		return r.Repo.Store[correlID], errs.ErrUniqueData
	}

	// Add record in Map.
	r.Repo.mu.Lock()
	defer r.Repo.mu.Unlock()

	r.Repo.Store[row.CorrelationID] = row
	r.Repo.uniqueFields[row.OriginalURL] = row.CorrelationID

	return row, nil
}
