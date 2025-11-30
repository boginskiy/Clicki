package repoMap

import (
	"context"

	"github.com/boginskiy/Clicki/internal/errs"
	"github.com/boginskiy/Clicki/internal/model"
)

type MapRecordsRepo struct {
	Repo *RepoMap
}

func NewMapRecordsRepo(repo *RepoMap) *MapRecordsRepo {
	return &MapRecordsRepo{
		Repo: repo,
	}
}

// CreateRecords.
func (r *MapRecordsRepo) CreateRecords(ctx context.Context, records any) error {
	rows, ok := records.([]model.ResURLSet)
	if !ok || len(rows) == 0 {
		return errs.NewErrPlace("records in CreateRecords not valid", nil)
	}
	r.Repo.mu.Lock()
	defer r.Repo.mu.Unlock()

	for _, row := range rows {
		// TODO! Перекладка с ResURLSet в URLTb не супер оптимально однако пока так ...
		r.Repo.Store[row.CorrelationID] = &model.URLTb{
			ID:            0,
			OriginalURL:   row.OriginalURL,
			ShortURL:      row.ShortURL,
			CorrelationID: row.CorrelationID,
			CreatedAt:     row.CreatedAt,
			UserID:        row.UserID,
		}
		r.Repo.uniqueFields[row.OriginalURL] = row.CorrelationID
	}
	return nil
}

// ReadRecords.
func (r *MapRecordsRepo) ReadRecords(ctx context.Context, userID int) (any, error) {
	records := []model.ResUserURLSet{}

	for _, v := range r.Repo.Store {
		if v.UserID == userID {
			records = append(records, model.ResUserURLSet{
				OriginalURL: v.OriginalURL,
				ShortURL:    v.ShortURL})
		}
	}
	return records, nil
}

// DeleteRecords.
func (r *MapRecordsRepo) DeleteRecords(ctx context.Context) error {
	return nil
}

// PingStore - for interface HealthCheckRepo.
func (r *MapRecordsRepo) PingStore(ctx context.Context) (bool, error) {
	return r.Repo.DB.CheckOpen()
}
