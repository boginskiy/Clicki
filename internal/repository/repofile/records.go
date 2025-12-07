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

func (fr *FileRecordsRepo) CreateRecords(ctx context.Context, records any) error {
	rows, ok := records.([]model.ResURLSet)
	if !ok || len(rows) == 0 {
		return errs.NewErrPlace("records in CreateRecords not valid", nil)
	}

	fr.Repo.mu.Lock()

	for _, row := range rows {
		fr.Repo.lastRecord += 1

		record := model.NewURLTb(fr.Repo.lastRecord, row.CorrelationID, row.OriginalURL, row.ShortURL, row.UserID)

		// Add data in Map.
		fr.Repo.tmpStore[record.CorrelationID] = record

		// Write data in file.
		jsonData, err := json.Marshal(record)
		if err != nil {
			return err
		}

		jsonData = append(jsonData, byte('\n'))
		_, err = fr.Repo.Store.Write(jsonData)
		if err != nil {
			return err
		}
	}
	fr.Repo.mu.Unlock()
	return nil
}

func (fr *FileRecordsRepo) ReadRecords(ctx context.Context, userID int) (any, error) {
	records := []model.ResUserURLSet{}

	fr.Repo.muR.RLock()
	defer fr.Repo.muR.RUnlock()

	for _, v := range fr.Repo.tmpStore {
		if v.UserID == userID {
			records = append(records, model.ResUserURLSet{
				OriginalURL: v.OriginalURL,
				ShortURL:    v.ShortURL})
		}
	}
	return records, nil
}

// DeleteRecords - for interface.
func (fr *FileRecordsRepo) DeleteRecords(ctx context.Context) error {
	return nil
}

// PingStore - for interface HealthCheckRepo.
func (fr *FileRecordsRepo) PingStore(ctx context.Context) (bool, error) {
	return fr.Repo.DB.CheckOpen()
}
