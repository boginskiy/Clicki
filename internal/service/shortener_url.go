package service

import (
	"github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/pkg"
)

const LONG = 8

var ShortenerURL = &ProURL{}

type ProURL struct {
}

func (b *ProURL) EncryptionLongURL() (imitationURL string) {
	for {
		// Вызов шифратора
		imitationURL = pkg.Scramble(LONG)
		// Проверка на уникальность
		if _, ok := db.Store[imitationURL]; !ok {
			break
		}
	}
	return imitationURL
}
