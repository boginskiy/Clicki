package service

import (
	"context"
	"net/http"
	"strings"

	c "github.com/boginskiy/Clicki/cmd/config"
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
	Kwargs      c.VarGetter
}

func NewShortURL(
	kwargs c.VarGetter, logger l.Logger, repo r.URLRepository,
	checker v.Checker, extraFuncer p.ExtraFuncer) *ShortURL {

	return &ShortURL{
		ExtraFuncer: extraFuncer,
		Checker:     checker,
		Logger:      logger,
		Kwargs:      kwargs,
		Repo:        repo,
	}
}

func (s *ShortURL) encryptionLongURL() (correlationID string) {
	for {
		correlationID = pkg.Scramble(LONG)                   // Вызов шифратора
		if s.Repo.CheckUnic(context.TODO(), correlationID) { // Проверка на уникальность
			break
		}
	}
	return correlationID
}

func (s *ShortURL) GetHeader() string {
	return "text/plain"
}

func (s *ShortURL) Create(req *http.Request) ([]byte, error) {
	originURL, err := s.ExtraFuncer.TakeAllBodyFromReq(req) // Вынимаем тело запроса

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

	correlationID := s.encryptionLongURL()                                      // Уникальный идентификатор
	shortURL := s.Kwargs.GetBaseURL() + "/" + correlationID                     // Новый сокращенный URL
	record := s.Repo.NewRow(context.TODO(), originURL, shortURL, correlationID) // Делаем запись для DB
	s.Repo.Create(context.TODO(), record)                                       // Кладем в DB данные

	return []byte(shortURL), nil
}

func (s *ShortURL) Read(req *http.Request) ([]byte, error) {
	correlationID := strings.TrimLeft(req.URL.Path, "/")      // Достаем параметр correlationID
	record, err := s.Repo.Read(context.TODO(), correlationID) // Достаем origin URL

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
	if s.Repo.GetDB() != nil {
		err := s.Repo.GetDB().Ping()
		if err != nil {
			return EmptyByteSlice, err
		}
	}
	return StoreDBIsSucces, nil
}

func (s *ShortURL) SetBatch(req *http.Request) ([]byte, error) {
	return StoreDBIsSucces, nil
}
