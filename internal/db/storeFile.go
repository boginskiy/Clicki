package db

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"sync"

	c "github.com/boginskiy/Clicki/cmd/config"
	l "github.com/boginskiy/Clicki/internal/logger"
	m "github.com/boginskiy/Clicki/internal/model"
)

const SIZE = 20

type StoreFile struct {
	Store   map[string]*m.URLFile
	scanner *bufio.Scanner
	mu      sync.RWMutex
	file    *os.File
	cntLine int
}

func NewStoreFile(kwargs c.VarGetter, _ l.Logger) (*StoreFile, error) {
	f, err := os.OpenFile(kwargs.GetPathToStore(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	sf := &StoreFile{file: f, scanner: bufio.NewScanner(f)}
	sf.Store = sf.dataRecovery()

	return sf, nil
}

func (sf *StoreFile) dataRecovery() map[string]*m.URLFile {
	result := make(map[string]*m.URLFile, SIZE)
	// Проход по строкам
	for sf.scanner.Scan() {
		record := &m.URLFile{}
		line := sf.scanner.Text()

		// Десериализация
		err := json.Unmarshal([]byte(line), record)
		if err != nil {
			continue
		}
		// Сохранение данных с map
		result[record.ShortURL] = record
		// Счетчик для UUID
		sf.cntLine = max(sf.cntLine, record.UUID)
	}
	return result
}

func (sf *StoreFile) GetDB() *sql.DB {
	return nil
}

func (sf *StoreFile) CloseDB() {
	sf.file.Close()
}

func (sf *StoreFile) CheckUnic(shortURL string) bool {
	_, ok := sf.Store[shortURL]
	return !ok
}

func (sf *StoreFile) NewRow(originURL, shortURL string) any {
	return m.NewURLFile(originURL, shortURL)
}

func (sf *StoreFile) Read(shortURL string) (any, error) {
	sf.mu.RLock()
	defer sf.mu.RUnlock()

	record, ok := sf.Store[shortURL]
	if !ok {
		return nil, errors.New("data is not available")
	}
	return record, nil
}

func (sf *StoreFile) Create(record any) error {
	row, ok := record.(*m.URLFile)
	if !ok {
		return errors.New("Error in StoreFile>Create")
	}

	sf.mu.Lock()
	sf.cntLine += 1
	row.UUID = sf.cntLine
	sf.Store[row.ShortURL] = row
	sf.mu.Unlock()

	tmpB, err := json.Marshal(row)
	if err != nil {
		return err
	}
	tmpB = append(tmpB, byte('\n'))
	_, err = sf.file.Write(tmpB)
	return err
}
