package tserv

import (
	"github.com/boginskiy/Clicki/internal/database"
	"github.com/boginskiy/Clicki/internal/model"
)

func UpdateDB(db database.DataBase) database.DataBase {
	switch v := db.GetDB().(type) {
	case map[string]*model.URLTb:

		record := &model.URLTb{
			ID:            0,
			OriginalURL:   "https://practicum.yandex.ru/",
			ShortURL:      "short_url",
			CorrelationID: "H3HIkks3",
			UserID:        100,
		}

		v["H3HIkks3"] = record
	}
	return db
}
