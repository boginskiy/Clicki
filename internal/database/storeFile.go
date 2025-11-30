package database

import (
	"os"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/errs"
	"github.com/boginskiy/Clicki/internal/logg"
)

// StoreFile - store is file database.
type StoreFile struct {
	Logg   logg.Logger
	File   *os.File
	isOpen bool
}

func NewStoreFile(config conf.Config, logger logg.Logger) (DataBase, error) {
	f, err := os.OpenFile(config.GetPathToStore(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &StoreFile{
		Logg:   logger,
		File:   f,
		isOpen: true,
	}, nil
}

func (sf *StoreFile) GetDB() any {
	return sf.File
}

func (sf *StoreFile) CloseDB() {
	sf.isOpen = false
	sf.File.Close()
}

func (sf *StoreFile) CheckOpen() (bool, error) {
	if !sf.isOpen {
		return false, errs.ErrPingDataBase
	}
	return sf.isOpen, nil
}
