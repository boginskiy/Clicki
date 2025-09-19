package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/db"
	e "github.com/boginskiy/Clicki/internal/errors"
	m "github.com/boginskiy/Clicki/internal/model"
	"github.com/jackc/pgerrcode"
)

type SQLURLRepository struct {
	Kwargs config.VarGetter
	DB     db.DBer
}

func NewSQLURLRepository(kwargs config.VarGetter, db db.DBer) *SQLURLRepository {
	return &SQLURLRepository{
		Kwargs: kwargs,
		DB:     db,
	}
}

func (s *SQLURLRepository) GetDB() *sql.DB {
	return s.DB.GetDB()
}

func (s *SQLURLRepository) CheckUnic(ctx context.Context, correlationID string) bool {
	return true
}

func (s *SQLURLRepository) Create(ctx context.Context, preRecord any) (any, error) {
	record, ok := preRecord.(*m.URLTb)
	if !ok {
		return nil, e.NewErrPlace("data not valid", nil)
	}

	errClassifier := NewPGErrorClass()
	tmpDB := s.DB.GetDB()

	// Strategy №2. SQl-Query-error.
	for attempt := 0; attempt <= s.Kwargs.GetMaxRetries(); attempt++ {

		row, err := InsertRowToUrls(tmpDB, context.TODO(),
			record.CorrelationID, record.OriginalURL, record.ShortURL, record.CreatedAt)

		// Ошибок нет, данные записаны
		if err == nil {
			id, _ := row.LastInsertId()
			record.ID = int(id)
			return record, nil
		}

		// Определяем поведение при получении ошибки
		errCode, needRetry := errClassifier.Classify(err)

		// Логика, если добавляемая запись не уникальна в БД
		if errCode == pgerrcode.UniqueViolation {

			// Делаем повторный запрос в БД
			row := SelectRowByOriginalURL(tmpDB, context.TODO(),
				record.OriginalURL)

			// Ошибок нет, возвращаем запись
			if err2 := row.Scan(
				&record.ID,
				&record.OriginalURL,
				&record.ShortURL,
				&record.CorrelationID,
				&record.CreatedAt); err2 == nil {

				return record, err
			} else {
				break
			}

			// Логика, если запрос к БД не надо повторять
		} else if needRetry == NonRetriable {
			break

			// Логика, если запрос к БД необходимо повторить
		} else {
			time.Sleep(3 * time.Millisecond)
		}
	}
	return nil, e.NewErrPlace("insert into is bad", nil)
}

func (s *SQLURLRepository) Read(ctx context.Context, correlID string) (any, error) {
	tmpDB := s.DB.GetDB()
	record := &m.URLTb{}

	row := SelectRowByCorrelID(tmpDB, context.TODO(), correlID)

	if err := row.Scan(
		&record.ID,
		&record.OriginalURL,
		&record.ShortURL,
		&record.CorrelationID,
		&record.CreatedAt); err != nil {
		return nil, err
	}

	return record, nil
}

func (s *SQLURLRepository) CreateSet(ctx context.Context, records any) error {
	rows, ok := records.([]m.ResURLSet)
	if !ok || len(rows) == 0 {
		return e.NewErrPlace("data not valid", nil)
	}

	tmpDB := s.DB.GetDB()

	tx, err := tmpDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, v := range rows {
		// все изменения записываются в транзакцию
		_, err := InsertRowToUrlsTX(tx, context.TODO(),
			v.CorrelationID, v.OriginalURL, v.ShortURL, v.CreatedAt)

		if err != nil {
			// если ошибка, то откатываем изменения
			tx.Rollback()
			return err
		}
	}
	// завершаем транзакцию
	tx.Commit()
	return nil
}
