package service

import (
	"log"
	"net/http"
	"strings"

	"github.com/boginskiy/Clicki/cmd/config"
	l "github.com/boginskiy/Clicki/internal/logger"
	m "github.com/boginskiy/Clicki/internal/model"
	p "github.com/boginskiy/Clicki/internal/preparation"
	r "github.com/boginskiy/Clicki/internal/repository"
	v "github.com/boginskiy/Clicki/internal/validation"
	"github.com/boginskiy/Clicki/pkg"
)

type ShortURL struct {
	ExtraFuncer p.ExtraFuncer
	Repo        r.URLRepository
	Checker     v.Checker
	Logger      l.Logger
}

func NewShortURL(
	repo r.URLRepository, logger l.Logger, checker v.Checker, extraFuncer p.ExtraFuncer) *ShortURL {
	log.Println(">>DB-3", repo.GetDB())
	return &ShortURL{
		ExtraFuncer: extraFuncer,
		Checker:     checker,
		Logger:      logger,
		Repo:        repo,
	}
}

func (s *ShortURL) encryptionLongURL() (shortURL string) {
	for {
		shortURL = pkg.Scramble(LONG)   // Вызов шифратора
		if s.Repo.CheckUnic(shortURL) { // Проверка на уникальность
			break
		}
	}
	return shortURL
}

func (s *ShortURL) Create(req *http.Request, kwargs config.VarGetter) ([]byte, error) {
	// Вынимаем тело запроса
	originURL, err := s.ExtraFuncer.TakeAllBodyFromReq(req)

	if err != nil {
		s.Logger.RaiseFatal(err, "ShortURL.Create>TakeAllBodyFromReq", nil)
		return EmptyByteSlice, err
	}

	// Валидируем URL. Проверка регуляркой, что строка является доменом сайта
	if !s.Checker.CheckUpURL(originURL) || originURL == "" {
		s.Logger.RaiseInfo("ShortURL.Create>CheckUpURL",
			l.Fields{"error": ErrDataNotValid.Error()})
		return EmptyByteSlice, ErrDataNotValid
	}

	shortURL := s.encryptionLongURL()            // Генерируем ключ
	record := s.Repo.NewRow(originURL, shortURL) // Делаем запись
	s.Repo.Create(record)                        // Кладем в db данные

	return []byte(shortURL), nil
}

func (s *ShortURL) Read(req *http.Request) ([]byte, error) {
	shortURL := strings.TrimLeft(req.URL.Path, "/") // Достаем параметр shortURL
	record, err := s.Repo.Read(shortURL)            // Достаем origin URL

	if err != nil {
		s.Logger.RaiseError(err, "ShortURL.Read>DB.Read", nil)
		return EmptyByteSlice, ErrDataNotValid
	}

	switch r := record.(type) {
	case *m.URLFile:
		return []byte(r.OriginalURL), nil
	case *m.URLTb:
		return []byte(r.OriginalURL), nil
	case string:
		return []byte(r), nil
	default:
		s.Logger.RaiseError(err, "ShortURL.Read>DB.Read>switch", nil)
		return EmptyByteSlice, ErrDataNotValid
	}
}

// CheckPing - check of connection db
func (s *ShortURL) CheckPing(req *http.Request) ([]byte, error) {

	err := s.Repo.GetDB().Ping()
	if err != nil {
		log.Println(">>4> Database connection is closed:", err)
	} else {
		log.Println(">>4> Database connection is active.")
	}

	// TODO! >>
	rows, err := s.Repo.GetDB().Query(
		`SELECT urls FROM information_schema.tables WHERE table_schema = 'public';`)
	if err != nil {
		log.Println(">>Err-4", err)

		rows, err = s.Repo.GetDB().Query(`SELECT * FROM urls;`)
		if err != nil {
			log.Println(">>Err-5", err)
			return EmptyByteSlice, err
		}

	}
	defer rows.Close()

	var tables []string
	log.Println("Err-6")
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			log.Fatalf("Scan error: %v\n", err)
		}
		tables = append(tables, tableName)
	}

	log.Println(">>Tables", tables)

	// Delete <<

	if s.Repo.GetDB() != nil {
		err := s.Repo.GetDB().Ping()
		if err != nil {
			return EmptyByteSlice, err
		}
	}
	return StoreDBIsSucces, nil
}
