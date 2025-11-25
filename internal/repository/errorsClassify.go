package repository

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

const (
	NonRetriable = 0 // NonRetriable - операцию не следует повторять.
	Retriable    = 1 // Retriable - операцию можно повторить.

)

// PGErrorClass - struct abount classification errors Database.
type PGErrorClass struct{}

func NewPGErrorClass() *PGErrorClass {
	return &PGErrorClass{}
}

// Classify - classification error and return code.
func (p *PGErrorClass) Classify(err error) (pq.ErrorCode, int) {
	if err == nil {
		return "", NonRetriable
	}

	// Check and convertation in pgconn.PgError, if it possible.
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return ClassifyPgError(pgErr)
	}

	// Default think that error doesn't repeate.
	return pgerrcode.Warning, NonRetriable
}

func ClassifyPgError(pgErr *pq.Error) (pq.ErrorCode, int) {
	switch pgErr.Code {

	// Grade 08 - Erorrs about connection.
	case pgerrcode.ConnectionException,
		pgerrcode.ConnectionDoesNotExist,
		pgerrcode.ConnectionFailure:
		return pgErr.Code, Retriable

	// Grade 40 - Erorrs about rollback of a transaction.
	case pgerrcode.TransactionRollback,
		pgerrcode.SerializationFailure,
		pgerrcode.DeadlockDetected:
		return pgErr.Code, Retriable

	// Grade 57 - Erorrs about operator.
	case pgerrcode.CannotConnectNow:
		return pgErr.Code, Retriable

	// Grade 22 - Erorrs about data.
	case pgerrcode.DataException,
		pgerrcode.NullValueNotAllowedDataException:
		return pgErr.Code, NonRetriable

	// Grade 23 - Erorrs about violation of integrity constraints.
	case pgerrcode.IntegrityConstraintViolation,
		pgerrcode.RestrictViolation,
		pgerrcode.NotNullViolation,
		pgerrcode.ForeignKeyViolation,
		pgerrcode.UniqueViolation,
		pgerrcode.CheckViolation:
		return pgErr.Code, NonRetriable

	// Grade 42 - Erorrs about syntactic.
	case pgerrcode.SyntaxErrorOrAccessRuleViolation,
		pgerrcode.SyntaxError,
		pgerrcode.UndefinedColumn,
		pgerrcode.UndefinedTable,
		pgerrcode.UndefinedFunction:
		return pgErr.Code, NonRetriable

	}
	return pgerrcode.Warning, NonRetriable
}
