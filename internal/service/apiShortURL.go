package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/boginskiy/Clicki/internal/logg"
	mod "github.com/boginskiy/Clicki/internal/model"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	rep "github.com/boginskiy/Clicki/internal/repository"
	valid "github.com/boginskiy/Clicki/internal/validation"
)

type APIShortURL struct {
	Repo        rep.Repository
	ExtraFuncer prep.ExtraFuncer
	Checker     valid.Checker
	Core        *CoreService
}

func NewAPIShortURL(
	core *CoreService,
	repo rep.Repository,
	checker valid.Checker,
	extraFuncer prep.ExtraFuncer) *APIShortURL {

	return &APIShortURL{
		ExtraFuncer: extraFuncer,
		Checker:     checker,
		Core:        core,
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
		s.Core.Logg.RaiseFatal(err, DeserializFatal, nil)
		return EmptyByteSlice, err
	}

	// Валидируем URL. Проверка регуляркой, что строка является доменом сайта
	if !s.Checker.CheckUpURL(bodyJSON.URL) || bodyJSON.URL == "" {
		s.Core.Logg.RaiseInfo("APIShortURL.Create>CheckUpURL",
			logg.Fields{"error": ErrDataNotValid.Error()})
		return EmptyByteSlice, ErrDataNotValid
	}

	userID := s.Core.TakeUserIDFromCtx(req)                      // Тащим идентификатор пользователя
	correlationID := s.Core.EncrypOriginURL()                    // Уникальный идентификатор
	shortURL := s.Core.Kwargs.GetBaseURL() + "/" + correlationID // Создаем новый сокращенный URL

	preRecord := mod.NewURLTb(0, correlationID, bodyJSON.URL, shortURL, userID) // Создаем черновую запись
	record, err := s.Repo.CreateRecord(context.TODO(), preRecord)               // Кладем в DB данные

	if err != nil && record == nil {
		s.Core.Logg.RaiseError(err, "APIShortURL.Create>Repo.Create", nil)
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
		s.Core.Logg.RaiseError(err, "APIShortURL.Create>switch", nil)
		return EmptyByteSlice, err
	}

	// Serialization
	result, err2 := s.ExtraFuncer.Serialization(resJSON)

	if err2 != nil {
		s.Core.Logg.RaiseError(err2, "APIShortURL.Create>NewResultJSON", nil)
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
		s.Core.Logg.RaiseFatal(ErrDataNotValid, "ShortURL>SetBatch>Token",
			logg.Fields{"fatal": ErrDataNotValid.Error()})
		return EmptyByteSlice, ErrDataNotValid
	}

	// Достаем идентификатор пользователя
	userID := s.Core.TakeUserIDFromCtx(req)

	// Разбор тела запроса
	respURLSet := make([]mod.ResURLSet, 0, 10)

	for decoder.More() {
		var rURL mod.ReqURLSet
		err := decoder.Decode(&rURL)

		if err != nil {
			s.Core.Logg.RaiseFatal(err, "ShortURL>SetBatch>Decode", nil)
			return EmptyByteSlice, err
		}

		shortURL := s.Core.Kwargs.GetBaseURL() + "/" + rURL.CorrelationID
		// Сбор множества URL
		respURLSet = append(respURLSet, mod.NewResURLSet(
			rURL.CorrelationID, rURL.OriginalURL, shortURL, userID))
	}

	// Сохранение в БД
	err := s.Repo.CreateRecords(context.TODO(), respURLSet)
	if err != nil {
		s.Core.Logg.RaiseError(err, "APIShortURL>SetBatch>CreateSet", nil)
		return EmptyByteSlice, err
	}

	// Сериализуем
	result, err := json.Marshal(respURLSet)
	s.Core.Logg.RaiseFatal(err, "ShortURL>SetBatch>Marshal", nil)
	return result, nil
}

func (s *APIShortURL) ReadSetUserURL(req *http.Request) ([]byte, error) {
	// Достаем идентификатор пользователя
	userID := s.Core.TakeUserIDFromCtx(req)

	dataSet, err := s.Repo.ReadRecords(context.TODO(), userID)
	if err != nil {
		s.Core.Logg.RaiseError(err, "APIShortURL.ReadSetUserURL>ReadSet", nil)
		return EmptyByteSlice, err
	}

	// Однозначно определяем наличие записей у пользователя
	records, ok := dataSet.([]mod.ResUserURLSet)
	if !ok {
		s.Core.Logg.RaiseError(ErrDataNotValid, "APIShortURL.ReadSetUserURL>Type?", nil)
		return EmptyByteSlice, ErrDataNotValid
	}
	if len(records) == 0 {
		return EmptyByteSlice, nil
	}

	// Serialization
	result, err := s.ExtraFuncer.Serialization(records)
	if err != nil {
		s.Core.Logg.RaiseError(err, "APIShortURL.ReadSetUserURL>Serialization", nil)
		return EmptyByteSlice, err
	}
	return result, err
}
