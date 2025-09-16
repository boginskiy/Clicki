package model

import (
	"fmt"
	"time"
)

type URLTb struct {
	id            int
	OriginalURL   string
	ShortURL      string
	CorrelationID string
	CreatedAt     time.Time
}

func NewURLTb(origin, short, correID string) *URLTb {
	return &URLTb{
		OriginalURL:   origin,
		ShortURL:      short,
		CorrelationID: correID,
		CreatedAt:     time.Now(),
	}
}

func (u *URLTb) String() string {
	return fmt.Sprintf(
		"id:%v, OriginalURL:%v, ShortURL:%v, CorrelationID:%v, CreatedAt:%v\n",
		u.id, u.OriginalURL, u.ShortURL, u.CorrelationID, u.CreatedAt)
}
