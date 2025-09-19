package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/boginskiy/Clicki/internal/db"
	e "github.com/boginskiy/Clicki/internal/errors"
	m "github.com/boginskiy/Clicki/internal/model"
	"github.com/jackc/pgerrcode"
)

type SQLURLRepository struct {
	DB db.DBer
}

func NewSQLURLRepository(db db.DBer) *SQLURLRepository {
	return &SQLURLRepository{
		DB: db,
	}
}

func (s *SQLURLRepository) convertTimeToStr(tm time.Time, pattern string) string {
	return tm.Format(pattern)
}

func (s *SQLURLRepository) convertStrToTime(tm string, pattern string) (time.Time, error) {
	return time.Parse(pattern, tm)
}

func (s *SQLURLRepository) CheckUnic(ctx context.Context, correlationID string) bool {
	return true
}

func (s *SQLURLRepository) GetDB() *sql.DB {
	return s.DB.GetDB()
}

func (s *SQLURLRepository) Create(ctx context.Context, preRecord any) (any, error) {
	record, ok := preRecord.(*m.URLTb)
	if !ok {
		return nil, e.NewErrPlace("data not valid", nil)
	}

	errClassifier := NewPGErrorClass()
	tmpDB := s.DB.GetDB()
	maxRetries := 3

	// Strategy №2. SQl-Query-error.
	for attempt := 0; attempt < maxRetries; attempt++ {

		row, err := tmpDB.ExecContext(context.TODO(),
			`INSERT INTO urls (correlation_id, original_url, short_url, created_at)
		 	 VALUES ($1, $2, $3, $4);`,
			record.CorrelationID, record.OriginalURL, record.ShortURL,
			s.convertTimeToStr(record.CreatedAt, time.RFC3339))

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
			row := tmpDB.QueryRowContext(context.TODO(),
				`SELECT id, original_url, short_url, correlation_id, created_at
				 FROM urls WHERE original_url = $1;`, record.OriginalURL)

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

func (s *SQLURLRepository) Read(ctx context.Context, correlationID string) (any, error) {
	tmpDB := s.DB.GetDB()
	tmpURL := &m.URLTb{}

	row := tmpDB.QueryRowContext(context.TODO(),
		`SELECT original_url, short_url, created_at 
		 FROM urls 
		 WHERE correlation_id = $1`,
		correlationID)

	var timeStr string

	err := row.Scan(&tmpURL.OriginalURL, &tmpURL.CorrelationID, &timeStr)
	if err != nil {
		return nil, err
	}

	timeT, err := s.convertStrToTime(timeStr, time.RFC3339)
	if err != nil {
		return nil, err
	}

	tmpURL.CreatedAt = timeT
	return tmpURL, nil
}

func (s *SQLURLRepository) Update(ctx context.Context, record *m.URLTb) {

}
func (s *SQLURLRepository) Delete(ctx context.Context, record *m.URLTb) {

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
		_, err := tx.ExecContext(ctx,
			`INSERT INTO urls (correlation_id, original_url, short_url, created_at)
		 	 VALUES($1,$2,$3,$4);`,
			v.CorrelationID, v.OriginalURL, v.ShortURL, s.convertTimeToStr(v.CreatedAt, time.RFC3339))

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

// TODO!
// Тестируем функционал всех видов БД
// Протянуть логгер, протянуть параметры ENV CLI | Количество ретрай
// Рефакторинг !!!!
// Set Batch
