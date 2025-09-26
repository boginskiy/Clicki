package service

import (
	"context"
	"encoding/json"
	"net/http"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
	mod "github.com/boginskiy/Clicki/internal/model"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/repository"
	valid "github.com/boginskiy/Clicki/internal/validation"
)

type APIShortURL struct {
	Repo        repository.Repository
	ExtraFuncer prep.ExtraFuncer
	Kwargs      conf.VarGetter
	Checker     valid.Checker
	Logger      logg.Logger
	Core        CoreServicer
}

func NewAPIShortURL(
	kwargs conf.VarGetter, logger logg.Logger, repo repository.Repository,
	checker valid.Checker, extraFuncer prep.ExtraFuncer) *APIShortURL {

	return &APIShortURL{
		Core:        NewCoreService(kwargs, logger, repo),
		ExtraFuncer: extraFuncer,
		Checker:     checker,
		Logger:      logger,
		Kwargs:      kwargs,
		Repo:        repo,
	}
}

func (s *APIShortURL) ReadURL(req *http.Request) ([]byte, error) {
	return EmptyByteSlice, nil
}

func (s *APIShortURL) CheckDB(req *http.Request) ([]byte, error) {
	return EmptyByteSlice, nil
}

func (s *APIShortURL) GetHeader() string {
	return "application/json"
}

func (s *APIShortURL) CreateURL(req *http.Request) ([]byte, error) {
	// Deserialization Body
	bodyJSON := mod.NewURLJson()
	err := s.ExtraFuncer.Deserialization(req, bodyJSON)

	if err != nil {
		s.Logger.RaiseFatal(err, DeserializFatal, nil)
		return EmptyByteSlice, err
	}

	// Валидируем URL. Проверка регуляркой, что строка является доменом сайта
	if !s.Checker.CheckUpURL(bodyJSON.URL) || bodyJSON.URL == "" {
		s.Logger.RaiseInfo("APIShortURL.CreateURL>CheckUpURL",
			logg.Fields{"error": ErrDataNotValid.Error()})
		return EmptyByteSlice, ErrDataNotValid
	}

	userID := s.Core.takeUserIdFromCtx(req)                 // Тащим идентификатор пользователя
	correlationID := s.Core.encrypOriginURL()               // Уникальный идентификатор
	shortURL := s.Kwargs.GetBaseURL() + "/" + correlationID // Создаем новый сокращенный URL

	preRecord := mod.NewURLTb(0, correlationID, bodyJSON.URL, shortURL, userID) // Создаем черновую запись
	record, err := s.Repo.Create(context.TODO(), preRecord)                     // Кладем в DB данные

	if err != nil && record == nil {
		s.Logger.RaiseError(err, "APIShortURL.CreateURL>Repo.Create", nil)
		return EmptyByteSlice, err
	}

	// Definition type
	var resJSON *mod.ResultJSON
	switch r := record.(type) {
	case *mod.URLTb:
		resJSON = mod.NewResultJSON(bodyJSON, r.ShortURL)
	case string:
		resJSON = mod.NewResultJSON(bodyJSON, r)
	default:
		s.Logger.RaiseError(err, "APIShortURL.CreateURL>switch", nil)
		return EmptyByteSlice, err
	}

	// Serialization Body
	result, err2 := s.ExtraFuncer.Serialization(resJSON)
	if err2 != nil {
		s.Logger.RaiseError(err2, "APIShortURL.CreateURL>NewResultJSON", nil)
		return EmptyByteSlice, err2
	}
	return result, err
}

func (s *APIShortURL) CreateSetURL(req *http.Request) ([]byte, error) {
	// Создаем декодер
	decoder := json.NewDecoder(req.Body)

	// Проверка, что пришло то, что надо
	token, _ := decoder.Token()
	if token != json.Delim('[') {
		s.Logger.RaiseFatal(ErrDataNotValid, "ShortURL>CreateSetURL>Token",
			logg.Fields{"fatal": ErrDataNotValid.Error()})
		return EmptyByteSlice, ErrDataNotValid
	}

	// Достаем идентификатор пользователя
	userID := s.Core.takeUserIdFromCtx(req)

	// Разбор тела запроса
	respURLSet := make([]mod.ResURLSet, 0, 10)

	for decoder.More() {
		var rURL mod.ReqURLSet
		err := decoder.Decode(&rURL)

		if err != nil {
			s.Logger.RaiseFatal(err, "ShortURL>CreateSetURL>Decode", nil)
			return EmptyByteSlice, err
		}

		shortURL := s.Kwargs.GetBaseURL() + "/" + rURL.CorrelationID

		// Сбор множества URL
		respURLSet = append(respURLSet, mod.NewResURLSet(
			rURL.CorrelationID, rURL.OriginalURL, shortURL, userID))
	}

	// Сохранение в БД
	err := s.Repo.CreateSet(context.TODO(), respURLSet)
	if err != nil {
		s.Logger.RaiseError(err, "APIShortURL>CreateSetURL>CreateSet", nil)
		return EmptyByteSlice, err
	}

	// Сериализуем
	result, err := json.Marshal(respURLSet)
	s.Logger.RaiseFatal(err, "ShortURL>CreateSetURL>Marshal", nil)
	return result, nil
}

func (s *APIShortURL) ReadSetUserURL(req *http.Request) ([]byte, error) {
	// Достаем идентификатор пользователя
	userID := s.Core.takeUserIdFromCtx(req)

	dataSet, err := s.Repo.ReadSet(context.TODO(), userID)
	if err != nil {
		s.Logger.RaiseError(err, "APIShortURL.ReadSetUserURL>ReadSet", nil)
		return EmptyByteSlice, err
	}

	// Однозначно определяем наличие записей у пользователя
	records, ok := dataSet.([]mod.ResUserURLSet)
	if !ok {
		s.Logger.RaiseError(ErrDataNotValid, "APIShortURL.ReadSetUserURL>Type?", nil)
		return EmptyByteSlice, ErrDataNotValid
	}
	if len(records) == 0 {
		return EmptyByteSlice, nil
	}

	// Serialization
	result, err := s.ExtraFuncer.Serialization(records)
	if err != nil {
		s.Logger.RaiseError(err, "APIShortURL.ReadSetUserURL>Serialization", nil)
		return EmptyByteSlice, err
	}
	return result, err
}
