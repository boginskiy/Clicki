package repofile

import (
	"context"
	"encoding/json"

	"github.com/boginskiy/Clicki/internal/errs"
	"github.com/boginskiy/Clicki/internal/model"
)

// generate:reset
type FileRecordsRepo struct {
	Repo      *RepoFile
	resetTest int
}

func NewFileRecordsRepo(repo *RepoFile) *FileRecordsRepo {
	return &FileRecordsRepo{
		Repo: repo,
	}
}

func (r *FileRecordsRepo) CreateRecords(ctx context.Context, records any) error {
	rows, ok := records.([]model.ResURLSet)
	if !ok || len(rows) == 0 {
		return errs.NewErrPlace("records in CreateRecords not valid", nil)
	}

	r.Repo.mu.Lock()

	for _, row := range rows {
		r.Repo.lastRecord += 1

		record := model.NewURLTb(r.Repo.lastRecord, row.CorrelationID, row.OriginalURL, row.ShortURL, row.UserID)

		// Add data in Map.
		r.Repo.tmpStore[record.CorrelationID] = record

		// Write data in file.
		jsonData, err := json.Marshal(record)
		if err != nil {
			return err
		}

		jsonData = append(jsonData, byte('\n'))
		_, err = r.Repo.Store.Write(jsonData)
		if err != nil {
			return err
		}
	}
	r.Repo.mu.Unlock()
	return nil
}

func (r *FileRecordsRepo) ReadRecords(ctx context.Context, userID int) (any, error) {
	records := []model.ResUserURLSet{}

	r.Repo.muR.RLock()
	defer r.Repo.muR.RUnlock()

	for _, v := range r.Repo.tmpStore {
		if v.UserID == userID {
			records = append(records, model.ResUserURLSet{
				OriginalURL: v.OriginalURL,
				ShortURL:    v.ShortURL})
		}
	}
	return records, nil
}

// DeleteRecords - for interface.
func (r *FileRecordsRepo) DeleteRecords(ctx context.Context) error {
	return nil
}

// PingStore - for interface HealthCheckRepo.
func (r *FileRecordsRepo) PingStore(ctx context.Context) (bool, error) {
	return r.Repo.DB.CheckOpen()
}
