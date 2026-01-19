package repodb

import (
	"context"
	"time"

	"github.com/boginskiy/Clicki/internal/errs"
	"github.com/boginskiy/Clicki/internal/model"
	"github.com/boginskiy/Clicki/internal/repository/utils"
	"github.com/jackc/pgerrcode"
)

type RepoDBRecord struct {
	Repo          *RepoDB
	errClassifier utils.DBErrClassifier
}

func NewRepoDBRecord(repoDB *RepoDB, errCl utils.DBErrClassifier) *RepoDBRecord {
	return &RepoDBRecord{
		Repo:          repoDB,
		errClassifier: errCl,
	}
}

func (r *RepoDBRecord) checkErrFromDB(ctx context.Context, record *model.URLTb, err error) (*model.URLTb, int) {
	code, needRetry := r.errClassifier.Classify(err)

	// Adding record not unique in DB.
	if code == pgerrcode.UniqueViolation {
		row := SelectRowByOriginalURL(r.Repo.Store, ctx,
			record.OriginalURL)

		// Ошибок нет, возвращаем запись.
		errScan := row.Scan(
			&record.ID,
			&record.OriginalURL,
			&record.ShortURL,
			&record.CorrelationID,
			&record.CreatedAt,
			&record.UserID)

		if errScan != nil {
			r.Repo.Logg.RaiseError(errScan, "error in repeated sending record from DB.", nil)
		}
		return record, needRetry
	}
	return nil, needRetry
}

// CreateRecord.
func (r *RepoDBRecord) CreateRecord(ctx context.Context, record *model.URLTb) (*model.URLTb, error) {
	for attempt := 0; attempt <= r.Repo.Cfg.GetMaxRetries(); attempt++ {
		row, errDB := InsertRowToUrls(ctx, r.Repo.Store, record)

		// There are not errors. Data is recorded.
		if errDB == nil {
			id, _ := row.LastInsertId()
			record.ID = int(id)
			return record, nil
		}

		// There are errors.
		record, needRetry := r.checkErrFromDB(ctx, record, errDB)
		if record != nil {
			// В ответ отдаю именно errDB для установки Conflict status.
			return record, errDB
		}

		if needRetry == utils.NonRetriable {
			// Logic, if query doesn't need to repeate for DB.
			break
		} else {
			// Logic, if query needs to repeate for DB.
			time.Sleep(10 * time.Millisecond)
		}
	}
	return nil, errs.NewErrPlace("bad insert to DB", nil)
}

// ReadRecord.
func (r *RepoDBRecord) ReadRecord(ctx context.Context, correlID string) (*model.URLTb, error) {
	record := &model.URLTb{}
	row := SelectRowByCorrelID(r.Repo.Store, ctx, correlID)

	if err := row.Scan(
		&record.ID,
		&record.OriginalURL,
		&record.ShortURL,
		&record.CorrelationID,
		&record.CreatedAt,
		&record.UserID,
		&record.DeletedFlag); err != nil {
		return nil, err
	}
	return record, nil
}

// CheckUniqueRecord - check that field correlationID about unique.
func (r *RepoDBRecord) CheckUniqueRecord(ctx context.Context, correlationID string) bool {
	// DataBase has unique field correlationID, therefore, always true.
	return true

}

// ReadLastRecord.
func (r *RepoDBRecord) ReadLastRecord(ctx context.Context) int {
	row := SelectMaxCntByUser(r.Repo.Store, ctx)
	var MaxCntByUser int

	err := row.Scan(&MaxCntByUser)
	if err != nil {
		return 0
	}
	return MaxCntByUser
}
