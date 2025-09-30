package repository

import (
	"context"
)

// Structures for communication of channels channels
type DelMessage struct {
	ListCorrelID []string
	UserID       int
}

func NewDelMessage(userID int) *DelMessage {
	return &DelMessage{UserID: userID}
}

type Repository interface {
	ReadRecord(ctx context.Context, recordID string) (any, error)
	CreateRecord(ctx context.Context, record any) (any, error)

	CheckUnicRecord(ctx context.Context, recordID string) bool
	CreateRecords(ctx context.Context, records any) error
	PingDB(ctx context.Context) (bool, error)

	// New
	ReadLastRecord(ctx context.Context) int
	ReadRecords(ctx context.Context, userID int) (any, error)
	MarkerRecords(ctx context.Context, messages ...DelMessage) error
	DeleteRecords(ctx context.Context) error
}
