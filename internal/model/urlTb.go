package model

import (
	"fmt"
	"time"
)

type URLTb struct {
	id          int
	OriginalURL string
	ShortURL    string
	CreatedAt   time.Time
}

func NewURLTb(origin, short string) *URLTb {
	return &URLTb{
		OriginalURL: origin,
		ShortURL:    short,
		CreatedAt:   time.Now(),
	}
}

func (u *URLTb) String() string {
	return fmt.Sprintf(
		"id:%v, OriginalURL:%v, ShortURL:%v, CreatedAt:%v\n",
		u.id, u.OriginalURL, u.ShortURL, u.CreatedAt)
}
