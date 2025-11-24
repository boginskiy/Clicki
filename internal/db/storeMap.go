package db

import (
	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/errs"
	"github.com/boginskiy/Clicki/internal/logg"
	mod "github.com/boginskiy/Clicki/internal/model"
)

const SIZE = 1024

type StoreMap struct {
	Store map[string]*mod.URLTb
}

func NewStoreMap(_ conf.Config, _ logg.Logger) (DBer, error) {
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
