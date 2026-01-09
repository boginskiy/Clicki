/*
RepoDBRecords - struct also implements interface: MarkerRepo, HealthCheckRepo.
*/
package repodb

import (
	"context"

	"github.com/boginskiy/Clicki/internal/errs"
	"github.com/boginskiy/Clicki/internal/model"
)

// RepoDBRecords - struct for processing records.
type RepoDBRecords struct {
	Repo *RepoDB
}

func NewRepoDBRecords(repo *RepoDB) *RepoDBRecords {
	return &RepoDBRecords{
		Repo: repo,
	}
}

// CreateRecords - implements interface RecordsRepo.
func (r *RepoDBRecords) CreateRecords(ctx context.Context, records any) error {
	rows, ok := records.([]model.ResURLSet)
	if !ok || len(rows) == 0 {
		return errs.NewErrPlace("records in CreateRecords not valid", nil)
	}

	tx, err := r.Repo.Store.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, v := range rows {
		// все изменения записываются в транзакцию.
		_, err := InsertRowToUrlsTX(tx, ctx,
			v.CorrelationID, v.OriginalURL, v.ShortURL, v.CreatedAt, v.UserID)

		if err != nil {
			// если ошибка, то откатываем изменения.
			tx.Rollback()
			return err
		}
	}
	// завершаем транзакцию.
	tx.Commit()
	return nil
}

// ReadRecords - implements interface RecordsRepo.
func (r *RepoDBRecords) ReadRecords(ctx context.Context, userID int) ([]model.ResUserURLSet, error) {
	records := []model.ResUserURLSet{}
	record := model.ResUserURLSet{}

	rows, err := SelectUserURLs(r.Repo.Store, ctx, userID)
	if err != nil {
		return nil, errs.NewErrPlace("data reading error in ReadRecords", err)
	}
	defer rows.Close()

	// Читаем данные.
	for rows.Next() {
		err := rows.Scan(&record.OriginalURL, &record.ShortURL)
		if err != nil {
			// TODO: Залогировать бы на всяк случай.
			continue
		}
		records = append(records, record)
	}
	if rows.Err() != nil {
		return nil, errs.NewErrPlace("scan is bad in ReadRecords", rows.Err())
	}
	return records, nil
}

// DeleteRecords.
func (r *RepoDBRecords) DeleteRecords(ctx context.Context) error {
	_, err := r.Repo.Store.ExecContext(ctx,
		`DELETE FROM urls
	 	 WHERE deleted_flag = TRUE;`)
	return err
}
