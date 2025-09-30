package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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

func (rd *RepositoryDBURL) CheckUnicRecord(ctx context.Context, correlationID string) bool {
	// TODO! Нужно натсроить DataBase
	// correlationID должно быть уникальное поле
	return true
}

func (rd *RepositoryDBURL) PingDB(ctx context.Context) (bool, error) {
	return rd.DB.CheckOpen()
}

func (rd *RepositoryDBURL) CreateRecord(ctx context.Context, preRecord any) (any, error) {
	record, ok := preRecord.(*mod.URLTb)
	if !ok {
		return nil, cerr.NewErrPlace("data not valid", nil)
	}

	errClassifier := NewPGErrorClass()

	// Strategy №2. SQl-Query-error.
	for attempt := 0; attempt <= rd.Kwargs.GetMaxRetries(); attempt++ {

		row, errDB := InsertRowToUrls(rd.db, ctx,
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
			row := SelectRowByOriginalURL(rd.db, ctx,
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

func (rd *RepositoryDBURL) ReadRecord(ctx context.Context, correlID string) (any, error) {
	record := &mod.URLTb{}
	row := SelectRowByCorrelID(rd.db, ctx, correlID)

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

func (rd *RepositoryDBURL) CreateRecords(ctx context.Context, records any) error {
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
		_, err := InsertRowToUrlsTX(tx, ctx,
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
func (rd *RepositoryDBURL) ReadLastRecord(ctx context.Context) int {
	row := SelectMaxCntByUser(rd.db, ctx)
	var MaxCntByUser int

	err := row.Scan(&MaxCntByUser)
	if err != nil {
		return 0
	}
	return MaxCntByUser
}

// New
func (rd *RepositoryDBURL) ReadRecords(ctx context.Context, userID int) (any, error) {
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
	if rows.Err() != nil {
		return nil, cerr.NewErrPlace("scan not good", rows.Err())
	}
	return records, nil
}

func (rd *RepositoryDBURL) MarkerRecords(ctx context.Context, messages ...DelMessage) error {
	// Для каждого пользователя создаем отдельный запрос
	// Мозгов хватило только на это. Реализации, когда все в одном запросе
	// привели меня к шизе.

	values := make([]string, 0, 10)
	args := make([]any, 0, 10)

	for _, mess := range messages {

		// Добавляем пользователя в аргументы
		args = append(args, mess.UserID) // $1

		for i, msg := range mess.ListCorrelID {
			values = append(values, fmt.Sprintf("$%d", (i+2)))
			args = append(args, msg)
		}

		query := `UPDATE urls
			SET deleted_flag = TRUE
			WHERE correlation_id IN (` + strings.Join(values, ",") + `)
			AND user_id = $1`

		_, err := rd.db.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
		// Обнуление перед следующей итерацией
		values = values[:0]
		args = args[:0]
	}
	return nil
}

func (rd *RepositoryDBURL) DeleteRecords(ctx context.Context) error {
	_, err := rd.db.ExecContext(ctx,
		`DELETE FROM urls
	 	 WHERE deleted_flag = TRUE;`)
	return err
}
