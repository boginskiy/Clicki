package tester

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/cmd/server"
	"github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/model"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/internal/validation"
)

func InitShortURL() *service.ShortURL {
	// Some part.
	logg := logg.NewLogg("test.log", "ERROR")
	config := config.NewVariables(logg)

	checker := validation.NewChecker()
	exFunc := prep.NewExtraFunc()

	// Db.
	layers := server.NewLayers(config, logg)
	db := layers.NewLayerDB()

	// add test data
	db = UpdateDB(db)

	repo := layers.NewLayerRepo(db)

	// Service.
	CoreServ := service.NewCoreService(config, logg, repo, nil)
	return service.NewShortURL(CoreServ, repo, checker, exFunc)
}

func UpdateDB(db db.DataBase) db.DataBase {
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
