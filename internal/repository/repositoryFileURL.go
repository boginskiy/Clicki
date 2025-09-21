package repository

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/db"
	cerr "github.com/boginskiy/Clicki/internal/error"
	mod "github.com/boginskiy/Clicki/internal/model"
)

const SIZE = 20

type RepositoryFileURL struct {
	Kwargs conf.VarGetter
	DB     db.DBer // *os.File

	Store        map[string]*mod.URLTb
	UniqueFields map[string]string
	Scanner      *bufio.Scanner
	muR          sync.RWMutex
	mu           sync.Mutex
	File         *os.File
	cntLine      int
}

func NewRepositoryFileURL(kwargs conf.VarGetter, dber db.DBer) (Repository, error) {
	// Create Scanner
	file, ok := dber.GetDB().(*os.File)
	if !ok {
		return nil, cerr.NewErrPlace("database not valid", nil)
	}
	// Create
	tmpRepo := &RepositoryFileURL{
		Kwargs:  kwargs,
		DB:      dber,
		Scanner: bufio.NewScanner(file),
		File:    file,
	}
	// Restoring data from the last session
	st, uf := tmpRepo.dataRecovery()
	tmpRepo.Store = st
	tmpRepo.UniqueFields = uf
	return tmpRepo, nil
}

func (rf *RepositoryFileURL) dataRecovery() (map[string]*mod.URLTb, map[string]string) {
	resultMap := make(map[string]*mod.URLTb, SIZE)
	resultSet := make(map[string]string, SIZE)

	// Проход по строкам
	for rf.Scanner.Scan() {
		record := &mod.URLTb{}
		line := rf.Scanner.Text()

		// Десериализация
		err := json.Unmarshal([]byte(line), record)
		if err != nil {
			continue
		}
		resultMap[record.CorrelationID] = record             // Сохранение данных с map
		resultSet[record.OriginalURL] = record.CorrelationID // Сохранение данных с set
		rf.cntLine = max(rf.cntLine, record.ID)              // Счетчик для UUID
	}
	return resultMap, resultSet
}

func (rf *RepositoryFileURL) CheckUnic(ctx context.Context, correlID string) bool {
	_, ok := rf.Store[correlID]
	return !ok
}

func (rf *RepositoryFileURL) Ping(ctx context.Context) (bool, error) {
	return rf.DB.CheckOpen()
}

func (sf *RepositoryFileURL) Read(ctx context.Context, correlID string) (any, error) {
	sf.muR.RLock()
	defer sf.muR.RUnlock()

	record, ok := sf.Store[correlID]
	if !ok {
		return nil, cerr.NewErrPlace("data is not available", nil)
	}
	return record, nil
}

func (sf *RepositoryFileURL) Create(ctx context.Context, preRecord any) (any, error) {
	row, ok := preRecord.(*mod.URLTb)
	if !ok {
		return nil, cerr.NewErrPlace("type is not available", nil)
	}

	// Логика, если данные уже есть в Store
	sf.muR.RLock()
	if correlID, ok := sf.UniqueFields[row.OriginalURL]; ok {
		return sf.Store[correlID], cerr.ErrUniqueData
	}
	sf.muR.RUnlock()

	// Логика, если данные отсутствуют в Store
	sf.mu.Lock()
	sf.cntLine += 1
	row.ID = sf.cntLine
	sf.Store[row.CorrelationID] = row
	sf.UniqueFields[row.OriginalURL] = row.CorrelationID
	sf.mu.Unlock()

	jsonData, err := json.Marshal(row)
	if err != nil {
		return nil, cerr.NewErrPlace("type is not available", err)
	}
	jsonData = append(jsonData, byte('\n'))

	// TODO! Может это вывести в интерфейс БД ? Может еще чего так же вывести и разгрузить репозиторий
	_, err = sf.File.Write(jsonData)
	return row, err
}

func (sf *RepositoryFileURL) CreateSet(ctx context.Context, records any) error {
	rows, ok := records.([]mod.ResURLSet)
	if !ok || len(rows) == 0 {
		return errors.New("data not valid")
	}

	sf.mu.Lock()

	for _, r := range rows {
		sf.cntLine += 1

		row := mod.NewURLTb(sf.cntLine, r.CorrelationID, r.OriginalURL, r.ShortURL)

		// Добавляем данные в Map
		sf.Store[row.CorrelationID] = row

		// Запись данных в файл
		jsonData, err := json.Marshal(row)
		if err != nil {
			return err
		}

		jsonData = append(jsonData, byte('\n'))
		_, err = sf.File.Write(jsonData)
		if err != nil {
			return err
		}
	}

	sf.mu.Unlock()
	return nil
}
