package repository

import (
	"context"
)

// Repository.
type Repository interface {
	HealthCheckRepo
	RecordsRepo
	RecordRepo
	MarkerRepo
}

// Record Storage.
type RecordRepo interface {
	CreateRecord(ctx context.Context, record any) (any, error)
	ReadRecord(ctx context.Context, recordID string) (any, error)
	CheckUniqueRecord(ctx context.Context, recordID string) bool
	ReadLastRecord(ctx context.Context) int
}

// Operations with groups of records.
type RecordsRepo interface {
	CreateRecords(ctx context.Context, records any) error
	ReadRecords(ctx context.Context, userID int) (any, error)
	DeleteRecords(ctx context.Context) error
}

// Marker entries.
type MarkerRepo interface {
	MarkRecords(ctx context.Context, messages any) error
}

// Diagnostic and monitoring methods.
type HealthCheckRepo interface {
	PingStore(ctx context.Context) (bool, error)
}
