package repodb

import (
	"context"
	"database/sql"
	"time"

	"github.com/boginskiy/Clicki/internal/repository/utils"
)

// InsertRowToUrls - add row in table 'urls'.
func InsertRowToUrls(db *sql.DB, ctx context.Context, id, origin, short string, tm time.Time, userID int) (sql.Result, error) {
	return db.ExecContext(ctx,
		`INSERT INTO urls (correlation_id, original_url, short_url, created_at, user_id)
		 VALUES ($1, $2, $3, $4, $5);`,
		id, origin, short, utils.ConvertTimeToStr(tm, time.RFC3339), userID)
}

// SelectRowByOriginalURL - choise row by field 'original_url'.
func SelectRowByOriginalURL(db *sql.DB, ctx context.Context, origin string) *sql.Row {
	return db.QueryRowContext(ctx,
		`SELECT id, original_url, short_url, correlation_id, created_at, user_id
		 FROM urls 
		 WHERE original_url = $1;`,
		origin)
}

// SelectRowByCorrelID - choise row by field 'correlation_id'.
func SelectRowByCorrelID(db *sql.DB, ctx context.Context, correlID string) *sql.Row {
	return db.QueryRowContext(ctx,
		`SELECT id, original_url, short_url, correlation_id, created_at, user_id, deleted_flag
		 FROM urls 
		 WHERE correlation_id = $1;`,
		correlID)
}

// InsertRowToUrlsTX - add row in table 'urls' through transaction.
func InsertRowToUrlsTX(tx *sql.Tx, ctx context.Context, id, origin, short string, tm time.Time, userID int) (sql.Result, error) {
	return tx.ExecContext(ctx,
		`INSERT INTO urls (correlation_id, original_url, short_url, created_at, user_id)
		 VALUES ($1, $2, $3, $4, $5);`,
		id, origin, short, utils.ConvertTimeToStr(tm, time.RFC3339), userID)
}

// SelectMaxCntByUser - .
func SelectMaxCntByUser(db *sql.DB, ctx context.Context) *sql.Row {
	return db.QueryRowContext(ctx,
		`SELECT MAX(user_id)
		 FROM urls`)
}

// IsThereUser - .
func IsThereUser(db *sql.DB, ctx context.Context, userID int) *sql.Row {
	return db.QueryRowContext(ctx,
		`SELECT EXISTS (
		 SELECT 1
		 FROM urls
		 WHERE user_id = $1)`,
		userID)
}

// SelectUserURLs - .
func SelectUserURLs(db *sql.DB, ctx context.Context, userID int) (*sql.Rows, error) {
	return db.QueryContext(ctx,
		`SELECT original_url, short_url 
		 FROM urls 
		 WHERE user_id = $1;`,
		userID)
}
