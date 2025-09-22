package repository

import (
	"context"
	"database/sql"
	"time"
)

// Query's *sql.DB

// InsertRowToUrls - добавление строки в таблицу 'urls'
func InsertRowToUrls(db *sql.DB, ctx context.Context, id, origin, short string, tm time.Time) (sql.Result, error) {
	return db.ExecContext(ctx,
		`INSERT INTO urls (correlation_id, original_url, short_url, created_at)
		 VALUES ($1, $2, $3, $4);`,
		id, origin, short, convertTimeToStr(tm, time.RFC3339))
}

// SelectRowByOriginalURL - выбор строки по полю 'original_url'
func SelectRowByOriginalURL(db *sql.DB, ctx context.Context, origin string) *sql.Row {
	return db.QueryRowContext(ctx,
		`SELECT id, original_url, short_url, correlation_id, created_at
		 FROM urls 
		 WHERE original_url = $1;`,
		origin)
}

// SelectRowByCorrelID - выбор строки по полю 'correlation_id'
func SelectRowByCorrelID(db *sql.DB, ctx context.Context, correlID string) *sql.Row {
	return db.QueryRowContext(context.TODO(),
		`SELECT id, original_url, short_url, correlation_id,  created_at
		 FROM urls 
		 WHERE correlation_id = $1;`,
		correlID)
}

// Query's *sql.Tx

// InsertRowToUrlsTX - добавление строки в таблицу 'urls' через транзакцию
func InsertRowToUrlsTX(tx *sql.Tx, ctx context.Context, id, origin, short string, tm time.Time) (sql.Result, error) {
	return tx.ExecContext(ctx,
		`INSERT INTO urls (correlation_id, original_url, short_url, created_at)
		 VALUES ($1, $2, $3, $4);`,
		id, origin, short, convertTimeToStr(tm, time.RFC3339))
}
