package server

import (
	"github.com/boginskiy/Clicki/internal/db"
	repo "github.com/boginskiy/Clicki/internal/repository"
)

type Designer interface {
	NewLayerRepo() repo.Repository
	NewLayerDB() db.DBer
	Close()
}
