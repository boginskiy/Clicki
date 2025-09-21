package db

import (
	"os"

	conf "github.com/boginskiy/Clicki/cmd/config"
	cerr "github.com/boginskiy/Clicki/internal/error"
	"github.com/boginskiy/Clicki/internal/logg"
)

type StoreFile struct {
	Logger logg.Logger
	File   *os.File
	isOpen bool
}

func NewStoreFile(kwargs conf.VarGetter, logger logg.Logger) (DBer, error) {
	f, err := os.OpenFile(kwargs.GetPathToStore(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &StoreFile{
		Logger: logger,
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
	if sf.isOpen == false {
		return false, cerr.ErrPingDataBase
	}
	return sf.isOpen, nil
}
