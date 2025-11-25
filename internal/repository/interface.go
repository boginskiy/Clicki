package repository

import (
	"context"
)

// Structures for communication of channels channels.
type DelMessage struct {
	ListCorrelID []string
	UserID       int64
}

func NewDelMessage(userID int64) *DelMessage {
	return &DelMessage{UserID: userID}
}

// Repository - .
type Repository interface {
	MarkerRecords(ctx context.Context, messages ...DelMessage) error
	ReadRecord(ctx context.Context, recordID string) (any, error)
	CreateRecord(ctx context.Context, record any) (any, error)
	CheckUnicRecord(ctx context.Context, recordID string) bool
	ReadRecords(ctx context.Context, userID int) (any, error)
	CreateRecords(ctx context.Context, records any) error
	PingDB(ctx context.Context) (bool, error)
	DeleteRecords(ctx context.Context) error
	ReadLastRecord(ctx context.Context) int
}
