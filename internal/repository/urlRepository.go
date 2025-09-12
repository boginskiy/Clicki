package repository

import (
	"context"
	"database/sql"
	"errors"
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

func (s *SQLURLRepository) NewRow(originURL, shortURL string) any {
	return m.NewURLTb(originURL, shortURL)
}

func (s *SQLURLRepository) CheckUnic(shortURL string) bool {
	return true
}

func (s *SQLURLRepository) GetDB() *sql.DB {
	return s.DB.GetDB()
}

func (s *SQLURLRepository) Create(record any) error {
	row, ok := record.(*m.URLTb)
	if !ok {
		return errors.New("error in SQLURLRepository>Create")
	}

	tmpDB := s.DB.GetDB()
	timeSrt := s.convertTimeToStr(row.CreatedAt, time.RFC3339)

	tmpDB.QueryRowContext(context.TODO(),
		`INSERT INTO urls (original_url, short_url, created_at)
		 VALUES ($1, $2, $3);`,
		row.OriginalURL, row.ShortURL, timeSrt)

	return nil
}

func (s *SQLURLRepository) Read(shortURL string) (any, error) {
	tmpDB := s.DB.GetDB()
	tmpURL := &m.URLTb{}

	row := tmpDB.QueryRowContext(context.TODO(),
		`SELECT original_url, short_url, created_at 
		 FROM urls 
		 WHERE short_url = $1`,
		shortURL)

	var timeStr string

	err := row.Scan(&tmpURL.OriginalURL, &tmpURL.ShortURL, &timeStr)
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

func (s *SQLURLRepository) Update(record *m.URLTb) {

}

func (s *SQLURLRepository) Delete(record *m.URLTb) {

}
