package repository

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/database"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/repository/repofile"
)

type MainRepoFile struct {
	HealthCheckRepo
	RecordsRepo
	RecordRepo
	MarkerRepo
}

func NewMainRepoFile(config config.Config, logger logg.Logger, db database.DataBase) Repository {
	// Init.
	RepoFile := repofile.NewRepoFile(config, logger, db)

	return &MainRepoFile{
		RecordsRepo:     repofile.NewFileRecordsRepo(RepoFile),
		RecordRepo:      repofile.NewFileRecordRepo(RepoFile),
		HealthCheckRepo: repofile.NewFileRecordsRepo(RepoFile),
		MarkerRepo:      nil, // repository is absent
	}
}
