package tserv

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/audit"
	"github.com/boginskiy/Clicki/internal/database"
	"github.com/boginskiy/Clicki/internal/layers"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/internal/tester/tfunc"
	"github.com/boginskiy/Clicki/internal/validation"
)

// InitAPIURLServ.
func InitAPIURLServ(logger logg.Logger, cfg config.Config) *service.APIURLServ {
	layers := layers.NewLayers(cfg, logger)
	db := layers.NewLayerDB()
	repo := layers.NewLayerRepo(db)

	tfunc.WriteRecord(repo)

	var sub1 = audit.NewFileReceiver(logger, cfg.GetAuditFile(), 1)
	var sub2 = audit.NewServerReceiver(logger, cfg.GetAuditURL(), 2)
	var publisher = audit.NewPublish(sub1, sub2)

	var fancer = preparation.NewFunctions(logger)
	var checker = validation.NewChecker()

	return service.NewAPIURLServ(cfg, logger, repo, checker, fancer, publisher)
}

// InitURLServ.
func InitURLServ(logger logg.Logger, cfg config.Config) *service.URLServ {
	db, _ := database.NewStoreMap(cfg, logger)
	repo := repository.NewMainRepoMap(cfg, logger, db)

	tfunc.WriteRecord(repo)

	var sub1 = audit.NewFileReceiver(logger, cfg.GetAuditFile(), 1)
	var sub2 = audit.NewServerReceiver(logger, cfg.GetAuditURL(), 2)
	var publisher = audit.NewPublish(sub1, sub2)

	var fancer = preparation.NewFunctions(logger)
	var checker = validation.NewChecker()

	return service.NewURLServ(cfg, logger, repo, checker, fancer, publisher)
}

func Init_URLServ_and_APIURLServ(logger logg.Logger, cfg config.Config) (*service.URLServ, *service.APIURLServ) {
	db, _ := database.NewStoreMap(cfg, logger)
	repo := repository.NewMainRepoMap(cfg, logger, db)

	tfunc.WriteRecord(repo)

	var sub1 = audit.NewFileReceiver(logger, cfg.GetAuditFile(), 1)
	var sub2 = audit.NewServerReceiver(logger, cfg.GetAuditURL(), 2)
	var publisher = audit.NewPublish(sub1, sub2)

	var fancer = preparation.NewFunctions(logger)
	var checker = validation.NewChecker()

	URLServ := service.NewURLServ(cfg, logger, repo, checker, fancer, publisher)
	APIURLServ := service.NewAPIURLServ(cfg, logger, repo, checker, fancer, publisher)

	return URLServ, APIURLServ
}
