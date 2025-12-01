package tester

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/audit"
	"github.com/boginskiy/Clicki/internal/layers"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/internal/validation"
)

func InitAPIURLServ(logger logg.Logger, cfg config.Config) *service.APIURLServ {
	layers := layers.NewLayers(cfg, logger)
	db := layers.NewLayerDB()
	repo := layers.NewLayerRepo(db)

	WriteRecord(repo)

	var sub1 = audit.NewFileReceiver(logger, cfg.GetAuditFile(), 1)
	var sub2 = audit.NewServerReceiver(logger, cfg.GetAuditURL(), 2)
	var publisher = audit.NewPublish(sub1, sub2)

	var fancer = preparation.NewFunctions()
	var checker = validation.NewChecker()

	return service.NewAPIURLServ(cfg, logger, repo, checker, fancer, publisher)
}
