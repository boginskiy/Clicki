package repository

import (
	"context"

	"github.com/boginskiy/Clicki/internal/model"
)

// Repository.
type Repository interface {
	StatisticianRepo
	HealthCheckRepo
	RecordsRepo
	RecordRepo
	MarkerRepo
}

// Record Storage.
type RecordRepo interface {
	CreateRecord(ctx context.Context, record *model.URLTb) (*model.URLTb, error)
	ReadRecord(ctx context.Context, recordID string) (*model.URLTb, error)
	CheckUniqueRecord(ctx context.Context, recordID string) bool
	ReadLastRecord(ctx context.Context) int
}

// Operations with groups of records.
type RecordsRepo interface {
	CreateRecords(ctx context.Context, records any) error
	ReadRecords(ctx context.Context, userID int) ([]model.ResUserURLSet, error)
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

// Collection stats.
type StatisticianRepo interface {
	ReadQuantityShortURLs(ctx context.Context) int
	ReadQuantityUsers(ctx context.Context) int
}
