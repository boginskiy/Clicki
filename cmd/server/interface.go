package server

import (
	"github.com/boginskiy/Clicki/internal/db"
	repo "github.com/boginskiy/Clicki/internal/repository"
)

// Designer - interface about manage Layers Appl.
type Designer interface {
	NewLayerRepo() repo.Repository
	NewLayerDB() db.DataBase
	Close()
}
