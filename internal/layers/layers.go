package layers

import (
	"errors"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/database"
	"github.com/boginskiy/Clicki/internal/logg"
	repo "github.com/boginskiy/Clicki/internal/repository"
)

// Designer - interface about manage Layers Appl.
type Designer interface {
	NewLayerRepo() repo.Repository
	NewLayerDB() database.DataBase
	Close()
}

// ErrCreateLayerDB - error about create layer DB.
var ErrCreateLayerDB = errors.New(`{"error":"layer of databse didn't create"}`)

// Layers - struct about layers of Appl.
type Layers struct {
	Cfg  conf.Config
	Logg logg.Logger
	db   database.DataBase
	rp   repo.Repository
}

func NewLayers(config conf.Config, logger logg.Logger) *Layers {
	return &Layers{
		Cfg:  config,
		Logg: logger,
	}
}

func (l *Layers) Close() {
	l.db.CloseDB()
}

// NewLayerDB - for create Repo layer.
func (l *Layers) NewLayerDB() database.DataBase {
	tmpdb, err := l.choiceLayerDB()
	if err != nil {
		l.Logg.RaiseFatal(err, "Layers>NewLayerDB", nil)
	}
	l.db = tmpdb
	return tmpdb
}

// NewLayerRepo - for create Repo layer.
func (l *Layers) NewLayerRepo(database database.DataBase) repo.Repository {
	if database == nil {
		l.Logg.RaiseFatal(ErrCreateLayerDB, "Layers>NewLayerRepo", nil)
	}
	l.rp = l.choiceLayerRepo(database)
	return l.rp
}

// ChoiceRepo - for a choise Repo layer.
func (l *Layers) choiceLayerRepo(db database.DataBase) repo.Repository {
	var newRepo func(conf.Config, logg.Logger, database.DataBase) repo.Repository

	switch db.(type) {
	case *database.StoreDB:
		newRepo = repo.NewMainRepoDB
	case *database.StoreFile:
		newRepo = repo.NewMainRepoFile
	case *database.StoreMap:
		newRepo = repo.NewMainRepoMap
	}
	return newRepo(l.Cfg, l.Logg, db)
}

// choiceLayerDB - for a choise DB layer.
func (l *Layers) choiceLayerDB() (database.DataBase, error) {
	var newStore func(conf.Config, logg.Logger) (database.DataBase, error)

	if l.Cfg.GetDB() != "" {
		newStore = database.NewStoreDB
	} else if l.Cfg.GetPathToStore() != "" {
		newStore = database.NewStoreFile
	} else {
		newStore = database.NewStoreMap
	}
	return newStore(l.Cfg, l.Logg)
}
