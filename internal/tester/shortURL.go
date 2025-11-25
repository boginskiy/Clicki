package tester

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/cmd/server"
	"github.com/boginskiy/Clicki/internal/logg"
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
	repo := layers.NewLayerRepo(db)

	// Service.
	CoreServ := service.NewCoreService(config, logg, repo, nil)
	return service.NewShortURL(CoreServ, repo, checker, exFunc)
}
