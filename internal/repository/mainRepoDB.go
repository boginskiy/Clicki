package repository

import (
	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/database"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/repository/repodb"
	"github.com/boginskiy/Clicki/internal/repository/utils"
)

type MainRepoDB struct {
	StatisticianRepo
	HealthCheckRepo
	RecordsRepo
	RecordRepo
	MarkerRepo
}

func NewMainRepoDB(config config.Config, logger logg.Logger, db database.DataBase) Repository {
	// Init
	RepoDB := repodb.NewRepoDB(config, logger, db)
	ErrClassifier := utils.NewPGErrorClass()

	return &MainRepoDB{
		HealthCheckRepo:  repodb.NewRepoDBHealthCheck(RepoDB),
		RecordsRepo:      repodb.NewRepoDBRecords(RepoDB),
		MarkerRepo:       repodb.NewRepoDBMarker(RepoDB),
		RecordRepo:       repodb.NewRepoDBRecord(RepoDB, ErrClassifier),
		StatisticianRepo: repodb.NewDBStatistician(RepoDB),
	}
}
