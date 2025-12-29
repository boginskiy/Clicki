package repository

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/database"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/repository/repomap"
)

type MainRepoMap struct {
	StatisticianRepo
	HealthCheckRepo
	RecordsRepo
	RecordRepo
	MarkerRepo
}

func NewMainRepoMap(config config.Config, logger logg.Logger, db database.DataBase) Repository {
	// Init.
	RepoMap := repomap.NewRepoMap(config, logger, db)

	return &MainRepoMap{
		RecordsRepo:      repomap.NewMapRecordsRepo(RepoMap),
		RecordRepo:       repomap.NewMapRecordRepo(RepoMap),
		HealthCheckRepo:  repomap.NewMapRecordsRepo(RepoMap),
		StatisticianRepo: repomap.NewMapStatistician(RepoMap),
		MarkerRepo:       nil, // repository is absent
	}
}
