package repodb

import (
	"database/sql"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/database"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/repository/utils"
)

// RepoDB - struct which Init one.
type RepoDB struct {
	Cfg   config.Config
	Logg  logg.Logger
	DB    database.DataBase
	Store *sql.DB
}

func NewRepoDB(config config.Config, logger logg.Logger, db database.DataBase) *RepoDB {
	store, ok := db.GetDB().(*sql.DB)
	if !ok {
		logger.RaiseFatal(utils.ErrInitRepoDB, "RepoDB", nil)
	}
	return &RepoDB{
		Cfg:   config,
		Logg:  logger,
		DB:    db,
		Store: store,
	}
}
