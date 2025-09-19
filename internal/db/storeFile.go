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
	e "github.com/boginskiy/Clicki/internal/errors"
	l "github.com/boginskiy/Clicki/internal/logger"
	m "github.com/boginskiy/Clicki/internal/model"
)

const SIZE = 20

type StoreFile struct {
	Store        map[string]*m.URLTb
	uniqueFields map[string]string
	scanner      *bufio.Scanner
	mu           sync.Mutex
	muR          sync.RWMutex
	file         *os.File
	cntLine      int
}

func NewStoreFile(kwargs c.VarGetter, _ l.Logger) (*StoreFile, error) {
	f, err := os.OpenFile(kwargs.GetPathToStore(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	sf := &StoreFile{
		file:    f,
		scanner: bufio.NewScanner(f),
	}
	sf.Store, sf.uniqueFields = sf.dataRecovery()
	return sf, nil
}

func (sf *StoreFile) dataRecovery() (map[string]*m.URLTb, map[string]string) {
	resultMap := make(map[string]*m.URLTb, SIZE)
	resultSet := make(map[string]string, SIZE)

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
		resultMap[record.CorrelationID] = record
		// Сохранение данных с set
		resultSet[record.OriginalURL] = record.CorrelationID
		// Счетчик для UUID
		sf.cntLine = max(sf.cntLine, record.ID)
	}
	return resultMap, resultSet
}

func (sf *StoreFile) GetDB() *sql.DB {
	return nil
}

func (sf *StoreFile) CloseDB() {
	sf.file.Close()
}

func (sf *StoreFile) CheckUnic(ctx context.Context, correlID string) bool {
	_, ok := sf.Store[correlID]
	return !ok
}

func (sf *StoreFile) Read(ctx context.Context, correlID string) (any, error) {
	sf.muR.RLock()
	defer sf.muR.RUnlock()

	record, ok := sf.Store[correlID]
	if !ok {
		return nil, e.NewErrPlace("data is not available", nil)
	}
	return record, nil
}

func (sf *StoreFile) Create(ctx context.Context, preRecord any) (any, error) {
	row, ok := preRecord.(*m.URLTb)
	if !ok {
		return nil, e.NewErrPlace("type is not available", nil)
	}

	// Логика, если данные уже есть в Store
	sf.muR.RLock()
	if correlID, ok := sf.uniqueFields[row.OriginalURL]; ok {
		return sf.Store[correlID], e.UniqueDataErr
	}
	sf.muR.RUnlock()

	// Логика, если данные отсутствуют в Store
	sf.mu.Lock()
	sf.cntLine += 1
	row.ID = sf.cntLine
	sf.Store[row.CorrelationID] = row
	sf.uniqueFields[row.OriginalURL] = row.CorrelationID
	sf.mu.Unlock()

	tmpB, err := json.Marshal(row)
	if err != nil {
		return nil, e.NewErrPlace("type is not available", err)
	}
	tmpB = append(tmpB, byte('\n'))
	_, err = sf.file.Write(tmpB)
	return row, err
}

func (sf *StoreFile) CreateSet(ctx context.Context, records any) error {
	rows, ok := records.([]m.ResURLSet)
	if !ok || len(rows) == 0 {
		return errors.New("data not valid")
	}

	sf.mu.Lock()

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

	sf.mu.Unlock()
	return nil
}
