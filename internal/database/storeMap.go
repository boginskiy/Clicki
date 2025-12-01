package database

import (
	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/errs"
	"github.com/boginskiy/Clicki/internal/logg"
	mod "github.com/boginskiy/Clicki/internal/model"
)

// SIZE is default size for Map.
const SIZE = 1024

// StoreMap - store is map database.
type StoreMap struct {
	Store map[string]*mod.URLTb
}

func NewStoreMap(_ conf.Config, _ logg.Logger) (DataBase, error) {
	return &StoreMap{
		Store: make(map[string]*mod.URLTb, SIZE),
	}, nil
}

func (sm *StoreMap) GetDB() any {
	return sm.Store
}

func (sm *StoreMap) CloseDB() {
}

func (sm *StoreMap) CheckOpen() (bool, error) {
	if sm.Store == nil {
		return false, errs.ErrPingDataBase
	}
	return true, nil
}
