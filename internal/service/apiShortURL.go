package service

import (
	"net/http"

	"github.com/boginskiy/Clicki/cmd/config"
	l "github.com/boginskiy/Clicki/internal/logger"
	m "github.com/boginskiy/Clicki/internal/model"
	p "github.com/boginskiy/Clicki/internal/preparation"
	r "github.com/boginskiy/Clicki/internal/repository"
	v "github.com/boginskiy/Clicki/internal/validation"
	"github.com/boginskiy/Clicki/pkg"
)

type APIShortURL struct {
	ExtraFuncer p.ExtraFuncer
	DB          r.URLRepository
	Checker     v.Checker
	Logger      l.Logger
}

func NewAPIShortURL(
	db r.URLRepository, logger l.Logger, checker v.Checker, extraFuncer p.ExtraFuncer) *APIShortURL {

	return &APIShortURL{
		ExtraFuncer: extraFuncer,
		Checker:     checker,
		Logger:      logger,
		DB:          db,
	}
}

func (s *APIShortURL) encryptionLongURL() (shortURL string) {
	for {
		shortURL = pkg.Scramble(LONG) // Вызов шифратора
		if s.DB.CheckUnic(shortURL) { // Проверка на уникальность
			break
		}
	}
	return shortURL
}

func (s *APIShortURL) Create(req *http.Request, kwargs config.VarGetter) ([]byte, error) {
	// Deserialization Body
	baseLink := m.NewURLJson()
	err := s.ExtraFuncer.Deserialization(req, baseLink)

	if err != nil {
		s.Logger.RaiseFatal(err, DeserializFatal, nil)
		return EmptyByteSlice, err
	}

	// Валидируем URL. Проверка регуляркой, что строка является доменом сайта
	if !s.Checker.CheckUpURL(baseLink.URL) || baseLink.URL == "" {
		s.Logger.RaiseInfo("APIShortURL.Create>CheckUpURL",
			l.Fields{"error": ErrDataNotValid.Error()})
		return EmptyByteSlice, ErrDataNotValid
	}

	shortURL := s.encryptionLongURL()            // Генерируем ключ
	record := m.NewURLTb(baseLink.URL, shortURL) // Создаем запись
	s.DB.Create(record)                          // Кладем в DB данные

	// Serialization Body
	extraLink := m.NewResultJSON(baseLink, kwargs.GetBaseURL()+"/"+shortURL)
	result, err := s.ExtraFuncer.Serialization(extraLink)

	if err != nil {
		s.Logger.RaiseError(err, "APIShortURL.Create>NewResultJSON", nil)
		return EmptyByteSlice, err
	}
	return result, nil
}

func (s *APIShortURL) Read(req *http.Request) ([]byte, error) {
	return EmptyByteSlice, nil
}

func (s *APIShortURL) CheckPing(req *http.Request) ([]byte, error) {
	return EmptyByteSlice, nil
}
