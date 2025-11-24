package server

import (
	"errors"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/logg"
	repo "github.com/boginskiy/Clicki/internal/repository"
)

var ErrCreateLayerDB = errors.New(`{"error":"layer of databse didn't create"}`)

type Layers struct {
	Cfg  conf.Config
	Logg logg.Logger
	db   db.DBer
	rp   repo.Repository
}

// NewLayers -
func NewLayers(config conf.Config, logger logg.Logger) *Layers {
	return &Layers{
		Cfg:  config,
		Logg: logger,
	}
}

// Close -
func (l *Layers) Close() {
	l.db.CloseDB()
}

// NewLayerDB -
func (l *Layers) NewLayerDB() db.DBer {
	dber, err := l.choiceLayerDB()
	if err != nil {
		l.Logg.RaiseFatal(err, "Layers>NewLayerDB", nil)
	}
	l.db = dber
	return dber
}

// NewLayerRepo -
func (l *Layers) NewLayerRepo() repo.Repository {
	if l.db == nil {
		l.Logg.RaiseFatal(ErrCreateLayerDB, "Layers>NewLayerRepo", nil)
	}
	repository, err := l.choiceLayerRepo(l.db)
	if err != nil {
		l.Logg.RaiseFatal(err, "Layers>NewLayerRepo", nil)
	}
	l.rp = repository
	return repository
}

// ChoiceRepo - for create Repo layer
func (l *Layers) choiceLayerRepo(database db.DBer) (repo.Repository, error) {
	var newRepo func(conf.Config, db.DBer) (repo.Repository, error)

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

// choiceLayerDB - for create DB layer
func (l *Layers) choiceLayerDB() (db.DBer, error) {
	var newStore func(conf.Config, logg.Logger) (db.DBer, error)

	if l.Cfg.GetDB() != "" {
		newStore = db.NewStoreDB
	} else if l.Cfg.GetPathToStore() != "" {
		newStore = db.NewStoreFile
	} else {
		newStore = db.NewStoreMap
	}
	return newStore(l.Cfg, l.Logg)
}
