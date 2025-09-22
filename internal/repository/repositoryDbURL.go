package repository

import (
	"context"
	"database/sql"
	"time"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/db"
	cerr "github.com/boginskiy/Clicki/internal/error"
	mod "github.com/boginskiy/Clicki/internal/model"
	"github.com/jackc/pgerrcode"
)

type RepositoryDBURL struct {
	Kwargs conf.VarGetter
	DB     db.DBer // *sql.DB
	db     *sql.DB
}

func NewRepositoryDBURL(kwargs conf.VarGetter, dber db.DBer) (Repository, error) {
	database, ok := dber.GetDB().(*sql.DB)
	if !ok {
		return nil, cerr.NewErrPlace("database not valid", nil)
	}
	return &RepositoryDBURL{
		Kwargs: kwargs,
		DB:     dber,
		db:     database,
	}, nil
}

func (rd *RepositoryDBURL) CheckUnic(ctx context.Context, correlationID string) bool {
	// TODO! Нужно натсроить DataBase
	// correlationID должно быть уникальное поле
	return true
}

func (rd *RepositoryDBURL) Ping(ctx context.Context) (bool, error) {
	return rd.DB.CheckOpen()
}

func (rd *RepositoryDBURL) Create(ctx context.Context, preRecord any) (any, error) {
	record, ok := preRecord.(*mod.URLTb)
	if !ok {
		return nil, cerr.NewErrPlace("data not valid", nil)
	}

	errClassifier := NewPGErrorClass()
	DB, ok := rd.DB.GetDB().(*sql.DB)
	if !ok {
		return nil, cerr.NewErrPlace("database not valid", nil)
	}

	// Strategy №2. SQl-Query-error.
	for attempt := 0; attempt <= rd.Kwargs.GetMaxRetries(); attempt++ {

		row, errDB := InsertRowToUrls(DB, context.TODO(),
			record.CorrelationID, record.OriginalURL, record.ShortURL, record.CreatedAt)

		// Ошибок нет, данные записаны
		if errDB == nil {
			id, _ := row.LastInsertId()
			record.ID = int(id)
			return record, nil
		}

		// Определяем поведение при получении ошибки
		code, needRetry := errClassifier.Classify(errDB)

		// Логика, если добавляемая запись не уникальна в БД
		if code == pgerrcode.UniqueViolation {

			// Делаем повторный запрос в БД
			row := SelectRowByOriginalURL(DB, context.TODO(),
				record.OriginalURL)

			// Ошибок нет, возвращаем запись
			if errScan := row.Scan(
				&record.ID,
				&record.OriginalURL,
				&record.ShortURL,
				&record.CorrelationID,
				&record.CreatedAt); errScan == nil {

				// В ответ отдаю именно errDB для установки статуса ответа
				return record, errDB
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
	return nil, cerr.NewErrPlace("insert into is bad", nil)
}

func (rd *RepositoryDBURL) Read(ctx context.Context, correlID string) (any, error) {
	DB, ok := rd.DB.GetDB().(*sql.DB)
	if !ok {
		return nil, cerr.NewErrPlace("database not valid", nil)
	}

	record := &mod.URLTb{}

	row := SelectRowByCorrelID(DB, context.TODO(), correlID)

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

func (rd *RepositoryDBURL) CreateSet(ctx context.Context, records any) error {
	rows, ok := records.([]mod.ResURLSet)
	if !ok || len(rows) == 0 {
		return cerr.NewErrPlace("data not valid", nil)
	}

	DB, ok := rd.DB.GetDB().(*sql.DB)
	if !ok {
		return cerr.NewErrPlace("database not valid", nil)
	}

	tx, err := DB.BeginTx(ctx, nil)
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
