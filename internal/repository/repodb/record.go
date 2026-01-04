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
	Repo *RepoDB
}

func NewRepoDBRecord(repoDB *RepoDB) *RepoDBRecord {
	return &RepoDBRecord{
		Repo: repoDB,
	}
}

func (r *RepoDBRecord) checkErrFromDB() {

}

// CreateRecord.
func (r *RepoDBRecord) CreateRecord(ctx context.Context, record *model.URLTb) (any, error) {
	errClassifier := utils.NewPGErrorClass()

	// Strategy №2. SQl-Query-error.
	for attempt := 0; attempt <= r.Repo.Cfg.GetMaxRetries(); attempt++ {
		row, errDB := InsertRowToUrls(ctx, r.Repo.Store, record)

		// There are not errors. Data is recorded.
		if errDB == nil {
			id, _ := row.LastInsertId()
			record.ID = int(id)
			return record, nil
		}

		// There are errors.
		code, needRetry := errClassifier.Classify(errDB)

		// Логика, если добавляемая запись не уникальна в БД.
		if code == pgerrcode.UniqueViolation {

			// Делаем повторный запрос в БД.
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

			// Ошибок нет, возвращаем запись.
			if errScan == nil {
				// В ответ отдаю именно errDB для установки статуса ответа.
				return record, errDB
			} else {
				break
			}

			// Логика, если запрос к БД не надо повторять.
		} else if needRetry == utils.NonRetriable {
			break
			// Логика, если запрос к БД необходимо повторить.
		} else {
			time.Sleep(3 * time.Millisecond)
		}
	}
	return nil, errs.NewErrPlace("insert into is bad", nil)
}

// ReadRecord.
func (r *RepoDBRecord) ReadRecord(ctx context.Context, correlID string) (any, error) {
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
