package testinit

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/audit"
	"github.com/boginskiy/Clicki/internal/database"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/model"
	"github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/internal/tester"
	"github.com/boginskiy/Clicki/internal/validation"
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

func InitURLServ(logger logg.Logger, cfg config.Config) *service.URLServ {
	db, _ := database.NewStoreMap(cfg, logger)
	repo := repository.NewMainRepoMap(cfg, logger, db)

	tester.WriteRecord(repo)

	var sub1 = audit.NewFileReceiver(logger, cfg.GetAuditFile(), 1)
	var sub2 = audit.NewServerReceiver(logger, cfg.GetAuditURL(), 2)
	var publisher = audit.NewPublish(sub1, sub2)

	var fancer = preparation.NewFunctions()
	var checker = validation.NewChecker()

	return service.NewURLServ(cfg, logger, repo, checker, fancer, publisher)
}

func InitConfig() *config.Variables {
	return &config.Variables{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080",
		ArgsCLI:       &config.ArgsCLI{},
		ArgsENV: &config.ArgsENV{
			SoftDeleteTime: 10,
			HardDeleteTime: 20,
		},
	}
}
