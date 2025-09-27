package repository

import (
	"context"
	"database/sql"
	"time"
)

// Query's *sql.DB

// InsertRowToUrls - добавление строки в таблицу 'urls'
func InsertRowToUrls(db *sql.DB, ctx context.Context, id, origin, short string, tm time.Time, userID int) (sql.Result, error) {
	return db.ExecContext(ctx,
		`INSERT INTO urls (correlation_id, original_url, short_url, created_at, user_id)
		 VALUES ($1, $2, $3, $4, $5);`,
		id, origin, short, convertTimeToStr(tm, time.RFC3339), userID)
}

// SelectRowByOriginalURL - выбор строки по полю 'original_url'
func SelectRowByOriginalURL(db *sql.DB, ctx context.Context, origin string) *sql.Row {
	return db.QueryRowContext(ctx,
		`SELECT id, original_url, short_url, correlation_id, created_at, user_id
		 FROM urls 
		 WHERE original_url = $1;`,
		origin)
}

// SelectRowByCorrelID - выбор строки по полю 'correlation_id'
func SelectRowByCorrelID(db *sql.DB, ctx context.Context, correlID string) *sql.Row {
	return db.QueryRowContext(ctx,
		`SELECT id, original_url, short_url, correlation_id, created_at, user_id
		 FROM urls 
		 WHERE correlation_id = $1;`,
		correlID)
}

// Query's *sql.Tx

// InsertRowToUrlsTX - добавление строки в таблицу 'urls' через транзакцию
func InsertRowToUrlsTX(tx *sql.Tx, ctx context.Context, id, origin, short string, tm time.Time, userID int) (sql.Result, error) {
	return tx.ExecContext(ctx,
		`INSERT INTO urls (correlation_id, original_url, short_url, created_at, user_id)
		 VALUES ($1, $2, $3, $4, $5);`,
		id, origin, short, convertTimeToStr(tm, time.RFC3339), userID)
}

// SelectMaxCntByUser -
func SelectMaxCntByUser(db *sql.DB, ctx context.Context) *sql.Row {
	return db.QueryRowContext(ctx,
		`SELECT MAX(user_id)
		 FROM urls`)
}

// IsThereUser -
func IsThereUser(db *sql.DB, ctx context.Context, userID int) *sql.Row {
	return db.QueryRowContext(ctx,
		`SELECT EXISTS (
		 SELECT 1
		 FROM urls
		 WHERE user_id = $1)`,
		userID)
}

// SelectUserURLs
func SelectUserURLs(db *sql.DB, ctx context.Context, userID int) (*sql.Rows, error) {
	return db.QueryContext(ctx,
		`SELECT original_url, short_url 
		 FROM urls 
		 WHERE user_id = $1;`,
		userID)
}
