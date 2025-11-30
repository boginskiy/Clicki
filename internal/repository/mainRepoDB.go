package repository

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/database"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/repository/repodb"
)

type MainRepoDB struct {
	HealthCheckRepo
	RecordsRepo
	RecordRepo
	MarkerRepo
}

func NewMainRepoDB(config config.Config, logger logg.Logger, db database.DataBase) Repository {
	// Init
	RepoDB := repodb.NewRepoDB(config, logger, db)

	return &MainRepoDB{
		HealthCheckRepo: repodb.NewRepoDBHealthCheck(RepoDB),
		RecordsRepo:     repodb.NewRepoDBRecords(RepoDB),
		MarkerRepo:      repodb.NewRepoDBMarker(RepoDB),
		RecordRepo:      repodb.NewRepoDBRecord(RepoDB),
	}
}
