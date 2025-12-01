package utils

import (
	"errors"
	"time"
)

// Errors.
var (
	ErrInitRepoDB   = errors.New(`{"error": "initialization RepoDB is bad"}`)
	ErrInitRepoFile = errors.New(`{"error": "initialization RepoFile is bad"}`)
	ErrInitRepoMap  = errors.New(`{"error": "initialization RepoMap is bad"}`)
)

// Function for convert time to string.
func ConvertTimeToStr(tm time.Time, pattern string) string {
	return tm.Format(pattern)
}

// DelMessage is struct for communication of channels.
type DelMessage struct {
	ListCorrelID []string
	UserID       int64
}

func NewDelMessage(userID int64) *DelMessage {
	return &DelMessage{UserID: userID}
}
