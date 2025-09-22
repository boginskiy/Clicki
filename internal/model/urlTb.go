package model

import (
	"fmt"
	"time"
)

type URLTb struct {
	ID            int       `json:"uuid"`
	OriginalURL   string    `json:"original_url"`
	ShortURL      string    `json:"short_url"`
	CorrelationID string    `json:"correlation_id"`
	CreatedAt     time.Time `json:"created_at"`
}

func NewURLTb(id int, correID, origin, short string) *URLTb {
	return &URLTb{
		ID:            id,
		OriginalURL:   origin,
		ShortURL:      short,
		CorrelationID: correID,
		CreatedAt:     time.Now(),
	}
}

func (u *URLTb) String() string {
	return fmt.Sprintf(
		"id:%v, OriginalURL:%v, ShortURL:%v, CorrelationID:%v, CreatedAt:%v\n",
		u.ID, u.OriginalURL, u.ShortURL, u.CorrelationID, u.CreatedAt)
}
