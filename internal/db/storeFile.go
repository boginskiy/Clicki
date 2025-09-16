package db

import (
	"bufio"
	"context"
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
	Store   map[string]*m.URLTb
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

func (sf *StoreFile) dataRecovery() map[string]*m.URLTb {
	result := make(map[string]*m.URLTb, SIZE)
	// Проход по строкам
	for sf.scanner.Scan() {
		record := &m.URLTb{}
		line := sf.scanner.Text()

		// Десериализация
		err := json.Unmarshal([]byte(line), record)
		if err != nil {
			continue
		}
		// Сохранение данных с map
		result[record.CorrelationID] = record
		// Счетчик для UUID
		sf.cntLine = max(sf.cntLine, record.ID)
	}
	return result
}

func (sf *StoreFile) GetDB() *sql.DB {
	return nil
}

func (sf *StoreFile) CloseDB() {
	sf.file.Close()
}

func (sf *StoreFile) CheckUnic(ctx context.Context, correlationID string) bool {
	_, ok := sf.Store[correlationID]
	return !ok
}

func (sf *StoreFile) Read(ctx context.Context, correlationID string) (any, error) {
	sf.mu.RLock()
	defer sf.mu.RUnlock()

	record, ok := sf.Store[correlationID]
	if !ok {
		return nil, errors.New("data is not available")
	}
	return record, nil
}

func (sf *StoreFile) Create(ctx context.Context, preRecord any) error {
	row, ok := preRecord.(*m.URLTb)
	if !ok {
		return errors.New("error in StoreFile>Create")
	}

	sf.mu.Lock()
	sf.cntLine += 1
	row.ID = sf.cntLine
	sf.Store[row.CorrelationID] = row
	sf.mu.Unlock()

	tmpB, err := json.Marshal(row)
	if err != nil {
		return err
	}
	tmpB = append(tmpB, byte('\n'))
	_, err = sf.file.Write(tmpB)
	return err
}

func (sf *StoreFile) CreateSet(ctx context.Context, records any) error {
	rows, ok := records.([]m.ResURLSet)
	if !ok || len(rows) == 0 {
		return errors.New("data not valid")
	}

	sf.mu.RLock()

	for _, r := range rows {
		sf.cntLine += 1

		row := m.NewURLTb(sf.cntLine, r.CorrelationID, r.OriginalURL, r.ShortURL)

		// Добавляем данные в Map
		sf.Store[row.CorrelationID] = row

		// Запись данных в файл
		tmpB, err := json.Marshal(row)
		if err != nil {
			return err
		}

		tmpB = append(tmpB, byte('\n'))
		_, err = sf.file.Write(tmpB)
		if err != nil {
			return err
		}
	}

	sf.mu.RUnlock()
	return nil
}
