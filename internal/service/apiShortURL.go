package service

import (
	"context"
	"encoding/json"
	"net/http"

	c "github.com/boginskiy/Clicki/cmd/config"
	l "github.com/boginskiy/Clicki/internal/logger"
	m "github.com/boginskiy/Clicki/internal/model"
	p "github.com/boginskiy/Clicki/internal/preparation"
	r "github.com/boginskiy/Clicki/internal/repository"
	v "github.com/boginskiy/Clicki/internal/validation"
	"github.com/boginskiy/Clicki/pkg"
)

type APIShortURL struct {
	Repo        r.URLRepository
	ExtraFuncer p.ExtraFuncer
	Kwargs      c.VarGetter
	Checker     v.Checker
	Logger      l.Logger
}

func NewAPIShortURL(
	kwargs c.VarGetter, logger l.Logger, repo r.URLRepository,
	checker v.Checker, extraFuncer p.ExtraFuncer) *APIShortURL {

	return &APIShortURL{
		ExtraFuncer: extraFuncer,
		Checker:     checker,
		Logger:      logger,
		Kwargs:      kwargs,
		Repo:        repo,
	}
}

func (s *APIShortURL) encryptionLongURL() (shortURL string) {
	for {
		shortURL = pkg.Scramble(LONG)                   // Вызов шифратора
		if s.Repo.CheckUnic(context.TODO(), shortURL) { // Проверка на уникальность
			break
		}
	}
	return shortURL
}

func (s *APIShortURL) GetHeader() string {
	return "application/json"
}

func (s *APIShortURL) Create(req *http.Request) ([]byte, error) {
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

	correlationID := s.encryptionLongURL()                                         // Уникальный идентификатор
	shortURL := s.Kwargs.GetBaseURL() + "/" + correlationID                        // Создаем новый сокращенный URL
	record := s.Repo.NewRow(context.TODO(), baseLink.URL, shortURL, correlationID) // Создаем запись
	s.Repo.Create(context.TODO(), record)                                          // Кладем в DB данные

	// Serialization Body
	resJSON := m.NewResultJSON(baseLink, shortURL)
	result, err := s.ExtraFuncer.Serialization(resJSON)

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

func (s *APIShortURL) SetBatch(req *http.Request) ([]byte, error) {
	// Создаем декодер
	decoder := json.NewDecoder(req.Body)

	// Проверка, что пришло то, что надо
	token, _ := decoder.Token()
	if token != json.Delim('[') {
		s.Logger.RaiseFatal(ErrDataNotValid, "ShortURL>SetBatch>Token",
			l.Fields{"fatal": ErrDataNotValid.Error()})
		return EmptyByteSlice, ErrDataNotValid
	}

	// Разбор тела запроса
	respURLSet := make([]m.ResURLSet, 0, 10)

	for decoder.More() {
		var rURL m.ReqURLSet
		err := decoder.Decode(&rURL)

		if err != nil {
			s.Logger.RaiseFatal(err, "ShortURL>SetBatch>Decode", nil)
			return EmptyByteSlice, err
		}

		// TODO!
		// shortURL := s.Kwargs.GetBaseURL() + "/" + s.encryptionLongURL()
		shortURL := s.Kwargs.GetBaseURL() + "/" + rURL.CorrelationID

		// Сбор множества URL
		respURLSet = append(respURLSet, m.NewResURLSet(rURL.CorrelationID, rURL.OriginalURL, shortURL))
	}

	// Сохранение в БД
	err := s.Repo.CreateSet(context.TODO(), respURLSet)
	if err != nil {
		s.Logger.RaiseError(err, "APIShortURL>SetBatch>CreateSet", nil)
		return EmptyByteSlice, err
	}

	// Сериализуем
	result, err := json.Marshal(respURLSet)
	s.Logger.RaiseFatal(err, "ShortURL>SetBatch>Marshal", nil)
	return result, nil
}
