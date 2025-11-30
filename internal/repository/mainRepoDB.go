package repository

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/database"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/repository/repoDB"
)

type MainRepoDB struct {
	HealthCheckRepo
	RecordsRepo
	RecordRepo
	MarkerRepo
}

func NewMainRepoDB(config config.Config, logger logg.Logger, db database.DataBase) Repository {
	// Init
	RepoDB := repoDB.NewRepoDB(config, logger, db)

	return &MainRepoDB{
		HealthCheckRepo: repoDB.NewRepoDBHealthCheck(RepoDB),
		RecordsRepo:     repoDB.NewRepoDBRecords(RepoDB),
		MarkerRepo:      repoDB.NewRepoDBMarker(RepoDB),
		RecordRepo:      repoDB.NewRepoDBRecord(RepoDB),
	}
}
