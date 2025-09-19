package repository

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

const (
	NonRetriable = 0 // NonRetriable - операцию не следует повторять
	Retriable    = 1 // Retriable - операцию можно повторить

)

type PGErrorClass struct{}

func NewPGErrorClass() *PGErrorClass {
	return &PGErrorClass{}
}

// Classify классифицирует ошибку и возвращает кодировку
func (p *PGErrorClass) Classify(err error) (pq.ErrorCode, int) {
	if err == nil {
		return "", NonRetriable
	}

	// Проверяем и конвертируем в pgconn.PgError, если это возможно
	// var pgErr *pgconn.PgError
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return ClassifyPgError(pgErr)
	}

	// По умолчанию считаем ошибку неповторяемой
	return pgerrcode.Warning, NonRetriable
}

func ClassifyPgError(pgErr *pq.Error) (pq.ErrorCode, int) {
	switch pgErr.Code {

	// Класс 08 - Ошибки соединения
	case pgerrcode.ConnectionException,
		pgerrcode.ConnectionDoesNotExist,
		pgerrcode.ConnectionFailure:
		return pgErr.Code, Retriable

	// Класс 40 - Откат транзакции
	case pgerrcode.TransactionRollback,
		pgerrcode.SerializationFailure,
		pgerrcode.DeadlockDetected:
		return pgErr.Code, Retriable

	// Класс 57 - Ошибка оператора
	case pgerrcode.CannotConnectNow:
		return pgErr.Code, Retriable

	// Класс 22 - Ошибки данных
	case pgerrcode.DataException,
		pgerrcode.NullValueNotAllowedDataException:
		return pgErr.Code, NonRetriable

	// Класс 23 - Нарушение ограничений целостности
	case pgerrcode.IntegrityConstraintViolation,
		pgerrcode.RestrictViolation,
		pgerrcode.NotNullViolation,
		pgerrcode.ForeignKeyViolation,
		pgerrcode.UniqueViolation,
		pgerrcode.CheckViolation:
		return pgErr.Code, NonRetriable

	// Класс 42 - Синтаксические ошибки
	case pgerrcode.SyntaxErrorOrAccessRuleViolation,
		pgerrcode.SyntaxError,
		pgerrcode.UndefinedColumn,
		pgerrcode.UndefinedTable,
		pgerrcode.UndefinedFunction:
		return pgErr.Code, NonRetriable

	}
	return pgerrcode.Warning, NonRetriable
}
