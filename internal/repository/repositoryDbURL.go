package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/errs"
	"github.com/boginskiy/Clicki/internal/model"
	"github.com/jackc/pgerrcode"
)

// RepositoryDBURL - repository for dataBase.
type RepositoryDBURL struct {
	Cfg conf.Config
	DB  db.DataBase
	db  *sql.DB
}

func NewRepositoryDBURL(config conf.Config, dataBase db.DataBase) (Repository, error) {
	tmpdb, ok := dataBase.GetDB().(*sql.DB)
	if !ok {
		return nil, errs.NewErrPlace("database not valid", nil)
	}
	return &RepositoryDBURL{
		Cfg: config,
		DB:  dataBase,
		db:  tmpdb,
	}, nil
}

// CheckUnicRecord - .
func (rd *RepositoryDBURL) CheckUnicRecord(ctx context.Context, correlationID string) bool {
	// TODO! Need settings DataBase. CorrelationID must be unic field.
	return true
}

// PingDB - .
func (rd *RepositoryDBURL) PingDB(ctx context.Context) (bool, error) {
	return rd.DB.CheckOpen()
}

// CreateRecord - .
func (rd *RepositoryDBURL) CreateRecord(ctx context.Context, preRecord any) (any, error) {
	record, ok := preRecord.(*model.URLTb)
	if !ok {
		return nil, errs.NewErrPlace("data not valid", nil)
	}

	errClassifier := NewPGErrorClass()

	// Strategy №2. SQl-Query-error.
	for attempt := 0; attempt <= rd.Cfg.GetMaxRetries(); attempt++ {

		row, errDB := InsertRowToUrls(rd.db, ctx,
			record.CorrelationID,
			record.OriginalURL,
			record.ShortURL,
			record.CreatedAt,
			record.UserID)

		// There are not errors. Data is recorded.
		if errDB == nil {
			id, _ := row.LastInsertId()
			record.ID = int(id)
			return record, nil
		}

		// Behaviour with gotting errors.
		code, needRetry := errClassifier.Classify(errDB)

		// Логика, если добавляемая запись не уникальна в БД.
		if code == pgerrcode.UniqueViolation {

			// Делаем повторный запрос в БД.
			row := SelectRowByOriginalURL(rd.db, ctx,
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
		} else if needRetry == NonRetriable {
			break
			// Логика, если запрос к БД необходимо повторить.
		} else {
			time.Sleep(3 * time.Millisecond)
		}
	}
	return nil, errs.NewErrPlace("insert into is bad", nil)
}

// ReadRecord - .
func (rd *RepositoryDBURL) ReadRecord(ctx context.Context, correlID string) (any, error) {
	record := &model.URLTb{}
	row := SelectRowByCorrelID(rd.db, ctx, correlID)

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

func (rd *RepositoryDBURL) CreateRecords(ctx context.Context, records any) error {
	rows, ok := records.([]model.ResURLSet)
	if !ok || len(rows) == 0 {
		return errs.NewErrPlace("data not valid", nil)
	}

	tx, err := rd.db.BeginTx(ctx, nil)
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

// ReadLastRecord - .
func (rd *RepositoryDBURL) ReadLastRecord(ctx context.Context) int {
	row := SelectMaxCntByUser(rd.db, ctx)
	var MaxCntByUser int

	err := row.Scan(&MaxCntByUser)
	if err != nil {
		return 0
	}
	return MaxCntByUser
}

// ReadRecords - .
func (rd *RepositoryDBURL) ReadRecords(ctx context.Context, userID int) (any, error) {
	records := []model.ResUserURLSet{}
	record := model.ResUserURLSet{}

	rows, err := SelectUserURLs(rd.db, ctx, userID)
	if err != nil {
		return nil, errs.NewErrPlace("data not valid", nil)
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
		return nil, errs.NewErrPlace("scan not good", rows.Err())
	}
	return records, nil
}

// MarkerRecords - .
func (rd *RepositoryDBURL) MarkerRecords(ctx context.Context, messages ...DelMessage) error {
	values := make([]string, 0, 10)
	args := make([]any, 0, 10)
	c := 1

	for _, mess := range messages {

		for _, correlID := range mess.ListCorrelID {
			values = append(values, fmt.Sprintf("($%d,$%d)", c, c+1))
			args = append(args, mess.UserID, correlID)
			c += 2
		}
	}

	query := fmt.Sprintf(`UPDATE urls
                          SET deleted_flag = TRUE
                          WHERE (user_id, correlation_id) IN (%s)`, strings.Join(values, ","))

	_, err := rd.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// DeleteRecords - .
func (rd *RepositoryDBURL) DeleteRecords(ctx context.Context) error {
	_, err := rd.db.ExecContext(ctx,
		`DELETE FROM urls
	 	 WHERE deleted_flag = TRUE;`)
	return err
}
