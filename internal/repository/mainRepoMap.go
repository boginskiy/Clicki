package repository

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/database"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/repository/repoMap"
)

type MainRepoMap struct {
	HealthCheckRepo
	RecordsRepo
	RecordRepo
	MarkerRepo
}

func NewMainRepoMap(config config.Config, logger logg.Logger, db database.DataBase) Repository {
	// Init.
	RepoMap := repoMap.NewRepoMap(config, logger, db)

	return &MainRepoMap{
		RecordsRepo:     repoMap.NewMapRecordsRepo(RepoMap),
		RecordRepo:      repoMap.NewMapRecordRepo(RepoMap),
		HealthCheckRepo: repoMap.NewMapRecordsRepo(RepoMap),
		MarkerRepo:      nil, // repository is absent
	}
}
