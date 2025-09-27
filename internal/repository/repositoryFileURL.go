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
	LastRec      int
	LastUser     int
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
		rf.LastRec = max(rf.LastRec, record.ID)              // Счетчик для ID
		rf.LastUser = max(rf.LastUser, record.UserID)        // Счетчик для User
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

func (rf *RepositoryFileURL) Read(ctx context.Context, correlID string) (any, error) {
	rf.muR.RLock()
	defer rf.muR.RUnlock()

	record, ok := rf.Store[correlID]
	if !ok {
		return nil, cerr.NewErrPlace("data is not available", nil)
	}
	return record, nil
}

func (rf *RepositoryFileURL) Create(ctx context.Context, preRecord any) (any, error) {
	row, ok := preRecord.(*mod.URLTb)
	if !ok {
		return nil, cerr.NewErrPlace("type is not available", nil)
	}

	// Логика, если данные уже есть в Store
	rf.muR.RLock()
	if correlID, ok := rf.UniqueFields[row.OriginalURL]; ok {
		return rf.Store[correlID], cerr.ErrUniqueData
	}
	rf.muR.RUnlock()

	// Логика, если данные отсутствуют в Store
	rf.mu.Lock()
	rf.LastRec += 1
	row.ID = rf.LastRec
	rf.Store[row.CorrelationID] = row
	rf.UniqueFields[row.OriginalURL] = row.CorrelationID
	rf.mu.Unlock()

	jsonData, err := json.Marshal(row)
	if err != nil {
		return nil, cerr.NewErrPlace("type is not available", err)
	}
	jsonData = append(jsonData, byte('\n'))

	// TODO! Может это вывести в интерфейс БД ? Может еще чего так же вывести и разгрузить репозиторий
	_, err = rf.File.Write(jsonData)
	return row, err
}

func (rf *RepositoryFileURL) CreateSet(ctx context.Context, records any) error {
	rows, ok := records.([]mod.ResURLSet)
	if !ok || len(rows) == 0 {
		return errors.New("data not valid")
	}

	rf.mu.Lock()

	for _, r := range rows {
		rf.LastRec += 1

		row := mod.NewURLTb(rf.LastRec, r.CorrelationID, r.OriginalURL, r.ShortURL, r.UserID)

		// Добавляем данные в Map
		rf.Store[row.CorrelationID] = row

		// Запись данных в файл
		jsonData, err := json.Marshal(row)
		if err != nil {
			return err
		}

		jsonData = append(jsonData, byte('\n'))
		_, err = rf.File.Write(jsonData)
		if err != nil {
			return err
		}
	}

	rf.mu.Unlock()
	return nil
}

// New
func (rf *RepositoryFileURL) TakeLastUser(ctx context.Context) int {
	return rf.LastUser
}

// New
func (rf *RepositoryFileURL) ReadSet(ctx context.Context, userID int) (any, error) {
	records := []mod.ResUserURLSet{}

	for _, v := range rf.Store {
		if v.UserID == userID {
			records = append(records, mod.ResUserURLSet{
				OriginalURL: v.OriginalURL,
				ShortURL:    v.ShortURL})
		}
	}
	return records, nil
}
