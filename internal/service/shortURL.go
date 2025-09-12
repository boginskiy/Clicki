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
	DB          r.URLRepository
	Checker     v.Checker
	Logger      l.Logger
}

func NewShortURL(
	db r.URLRepository, logger l.Logger, checker v.Checker, extraFuncer p.ExtraFuncer) *ShortURL {

	return &ShortURL{
		ExtraFuncer: extraFuncer,
		Checker:     checker,
		Logger:      logger,
		DB:          db,
	}
}

func (s *ShortURL) encryptionLongURL() (shortURL string) {
	for {
		shortURL = pkg.Scramble(LONG) // Вызов шифратора
		if s.DB.CheckUnic(shortURL) { // Проверка на уникальность
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

	shortURL := s.encryptionLongURL()          // Генерируем ключ
	record := s.DB.NewRow(originURL, shortURL) // Делаем запись
	s.DB.Create(record)                        // Кладем в db данные

	return []byte(shortURL), nil
}

func (s *ShortURL) Read(req *http.Request) ([]byte, error) {
	shortURL := strings.TrimLeft(req.URL.Path, "/") // Достаем параметр shortURL
	record, err := s.DB.Read(shortURL)              // Достаем origin URL

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

	// TODO! >>
	rows, err := s.DB.GetDB().Query(
		`SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';`)
	if err != nil {
		log.Println("Err", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			log.Fatalf("Scan error: %v\n", err)
		}
		tables = append(tables, tableName)
	}

	log.Println("Tables", tables)

	// Delete <<

	if s.DB.GetDB() != nil {
		err := s.DB.GetDB().Ping()
		if err != nil {
			return EmptyByteSlice, err
		}
	}
	return StoreDBIsSucces, nil
}
