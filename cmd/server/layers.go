package server

import (
	"errors"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/logg"
	repo "github.com/boginskiy/Clicki/internal/repository"
)

// ErrCreateLayerDB - error about create layer DB.
var ErrCreateLayerDB = errors.New(`{"error":"layer of databse didn't create"}`)

// Layers - struct about layers of Appl.
type Layers struct {
	Cfg  conf.Config
	Logg logg.Logger
	db   db.DataBase
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
func (l *Layers) NewLayerDB() db.DataBase {
	tmpdb, err := l.choiceLayerDB()
	if err != nil {
		l.Logg.RaiseFatal(err, "Layers>NewLayerDB", nil)
	}
	l.db = tmpdb
	return tmpdb
}

// NewLayerRepo - for create Repo layer.
func (l *Layers) NewLayerRepo(database db.DataBase) repo.Repository {
	if database == nil {
		l.Logg.RaiseFatal(ErrCreateLayerDB, "Layers>NewLayerRepo", nil)
	}
	repository, err := l.choiceLayerRepo(database)
	if err != nil {
		l.Logg.RaiseFatal(err, "Layers>NewLayerRepo", nil)
	}
	l.rp = repository
	return repository
}

// ChoiceRepo - for a choise Repo layer.
func (l *Layers) choiceLayerRepo(database db.DataBase) (repo.Repository, error) {
	var newRepo func(conf.Config, db.DataBase) (repo.Repository, error)

	switch database.(type) {
	case *db.StoreDB:
		newRepo = repo.NewRepositoryDBURL
	case *db.StoreFile:
		newRepo = repo.NewRepositoryFileURL
	case *db.StoreMap:
		newRepo = repo.NewRepositoryMapURL
	}
	return newRepo(l.Cfg, database)
}

// choiceLayerDB - for a choise DB layer.
func (l *Layers) choiceLayerDB() (db.DataBase, error) {
	var newStore func(conf.Config, logg.Logger) (db.DataBase, error)

	if l.Cfg.GetDB() != "" {
		newStore = db.NewStoreDB
	} else if l.Cfg.GetPathToStore() != "" {
		newStore = db.NewStoreFile
	} else {
		newStore = db.NewStoreMap
	}
	return newStore(l.Cfg, l.Logg)
}
