package repoFile

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/database"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/model"
	"github.com/boginskiy/Clicki/internal/repository/utils"
)

// RepoFile - repository for file.
type RepoFile struct {
	Cfg   config.Config
	Logg  logg.Logger
	DB    database.DataBase
	Store *os.File

	tmpStore     map[string]*model.URLTb
	uniqueFields map[string]string
	lastRecord   int
	lastUser     int

	scanner *bufio.Scanner
	muR     sync.RWMutex
	mu      sync.Mutex
}

func NewRepoFile(config config.Config, logger logg.Logger, db database.DataBase) *RepoFile {
	store, ok := db.GetDB().(*os.File)
	if !ok {
		logger.RaiseFatal(utils.ErrInitRepoDB, "RepoFile", nil)
	}

	// Create.
	tmpRf := &RepoFile{
		Cfg:     config,
		DB:      db,
		scanner: bufio.NewScanner(store),
		Store:   store,
	}
	// Restoring data from the last session.
	tmpRf.tmpStore, tmpRf.uniqueFields = tmpRf.dataRecovery()
	return tmpRf
}

// dataRecovery - .
func (r *RepoFile) dataRecovery() (map[string]*model.URLTb, map[string]string) {
	resultMap := make(map[string]*model.URLTb, database.SIZE)
	resultSet := make(map[string]string, database.SIZE)

	// Проход по строкам.
	for r.scanner.Scan() {
		record := &model.URLTb{}
		line := r.scanner.Text()

		// Deserialization.
		err := json.Unmarshal([]byte(line), record)
		if err != nil {
			continue
		}
		resultMap[record.CorrelationID] = record             // Data save with Map.
		resultSet[record.OriginalURL] = record.CorrelationID // Data save with Set.
		r.lastRecord = max(r.lastRecord, record.ID)          // Counter ID.
		r.lastUser = max(r.lastUser, record.UserID)          // Counter User.
	}
	return resultMap, resultSet
}
