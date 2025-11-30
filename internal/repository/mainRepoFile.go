package repository

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/database"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/repository/repoFile"
)

type MainRepoFile struct {
	HealthCheckRepo
	RecordsRepo
	RecordRepo
	MarkerRepo
}

func NewMainRepoFile(config config.Config, logger logg.Logger, db database.DataBase) Repository {
	// Init.
	RepoFile := repoFile.NewRepoFile(config, logger, db)

	return &MainRepoFile{
		RecordsRepo:     repoFile.NewFileRecordsRepo(RepoFile),
		RecordRepo:      repoFile.NewFileRecordRepo(RepoFile),
		HealthCheckRepo: repoFile.NewFileRecordsRepo(RepoFile),
		MarkerRepo:      nil, // repository is absent
	}
}
