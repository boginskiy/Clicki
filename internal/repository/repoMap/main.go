package repoMap

import (
	"sync"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/database"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/model"
	"github.com/boginskiy/Clicki/internal/repository/utils"
)

type RepoMap struct {
	Cfg   config.Config
	Logg  logg.Logger
	DB    database.DataBase
	Store map[string]*model.URLTb

	uniqueFields map[string]string
	muR          sync.RWMutex
	mu           sync.Mutex
}

func NewRepoMap(config config.Config, logger logg.Logger, db database.DataBase) *RepoMap {
	store, ok := db.GetDB().(map[string]*model.URLTb)
	if !ok {
		logger.RaiseFatal(utils.ErrInitRepoMap, "RepoMap", nil)
	}

	return &RepoMap{
		Cfg:          config,
		DB:           db,
		uniqueFields: make(map[string]string, database.SIZE),
		Store:        store,
	}
}
