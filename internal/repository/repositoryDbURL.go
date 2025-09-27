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

	// Strategy №2. SQl-Query-error.
	for attempt := 0; attempt <= rd.Kwargs.GetMaxRetries(); attempt++ {

		row, errDB := InsertRowToUrls(rd.db, context.TODO(),
			record.CorrelationID,
			record.OriginalURL,
			record.ShortURL,
			record.CreatedAt,
			record.UserID)

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
			row := SelectRowByOriginalURL(rd.db, context.TODO(),
				record.OriginalURL)

			// Ошибок нет, возвращаем запись
			errScan := row.Scan(
				&record.ID,
				&record.OriginalURL,
				&record.ShortURL,
				&record.CorrelationID,
				&record.CreatedAt,
				&record.UserID)

			// Ошибок нет, возвращаем запись
			if errScan == nil {
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
	record := &mod.URLTb{}
	row := SelectRowByCorrelID(rd.db, context.TODO(), correlID)

	if err := row.Scan(
		&record.ID,
		&record.OriginalURL,
		&record.ShortURL,
		&record.CorrelationID,
		&record.CreatedAt,
		&record.UserID); err != nil {
		return nil, err
	}
	return record, nil
}

func (rd *RepositoryDBURL) CreateSet(ctx context.Context, records any) error {
	rows, ok := records.([]mod.ResURLSet)
	if !ok || len(rows) == 0 {
		return cerr.NewErrPlace("data not valid", nil)
	}

	tx, err := rd.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, v := range rows {
		// все изменения записываются в транзакцию
		_, err := InsertRowToUrlsTX(tx, context.TODO(),
			v.CorrelationID, v.OriginalURL, v.ShortURL, v.CreatedAt, v.UserID)

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

// New
func (rd *RepositoryDBURL) TakeLastUser(ctx context.Context) (int, error) {
	row := SelectMaxCntByUser(rd.db, context.TODO())
	var MaxCntByUser int

	err := row.Scan(&MaxCntByUser)
	if err != nil {
		return -1, cerr.NewErrPlace("scan row is bad", err)
	}
	return MaxCntByUser, nil
}

// New
func (rd *RepositoryDBURL) CheckUser(ctx context.Context, userID int) (bool, error) {
	var exists bool
	row := IsThereUser(rd.db, ctx, userID)
	err := row.Scan(&exists)
	if err != nil {
		return false, cerr.NewErrPlace("scan row is bad", err)
	}
	return exists, nil
}

// New
func (rd *RepositoryDBURL) ReadSet(ctx context.Context, userID int) (any, error) {
	records := []mod.ResUserURLSet{}
	record := mod.ResUserURLSet{}

	rows, err := SelectUserURLs(rd.db, ctx, userID)
	if err != nil {
		return nil, cerr.NewErrPlace("data not valid", nil)
	}
	defer rows.Close()

	// Читаем данные
	for rows.Next() {
		err := rows.Scan(&record.OriginalURL, &record.ShortURL)
		if err != nil {
			// TODO! Залогировать бы на всяк случай
			continue
		}
		records = append(records, record)
	}
	return records, nil
}
