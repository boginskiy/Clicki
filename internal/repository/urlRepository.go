package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/boginskiy/Clicki/internal/db"
	m "github.com/boginskiy/Clicki/internal/model"
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

func (s *SQLURLRepository) NewRow(ctx context.Context, originURL, shortURL, correlationID string) any {
	return m.NewURLTb(originURL, shortURL, correlationID)
}

func (s *SQLURLRepository) CheckUnic(ctx context.Context, correlation_id string) bool {
	return true
}

func (s *SQLURLRepository) GetDB() *sql.DB {
	return s.DB.GetDB()
}

func (s *SQLURLRepository) Create(ctx context.Context, record any) error {
	row, ok := record.(*m.URLTb)
	if !ok {
		return errors.New("error in SQLURLRepository>Create")
	}

	tmpDB := s.DB.GetDB()

	tmpDB.QueryRowContext(context.TODO(),
		`INSERT INTO urls (correlation_id, original_url, short_url, created_at)
		 VALUES ($1, $2, $3, $4);`,
		row.CorrelationID, row.OriginalURL, row.ShortURL, s.convertTimeToStr(row.CreatedAt, time.RFC3339))

	return nil
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
		return errors.New("data not valid")
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
			fmt.Println(err)
			tx.Rollback()
			return err
		}
	}
	// завершаем транзакцию
	tx.Commit()
	return nil
}
