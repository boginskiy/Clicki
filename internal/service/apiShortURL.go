package service

import (
	"context"
	"encoding/json"
	"net/http"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
	mod "github.com/boginskiy/Clicki/internal/model"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	repo "github.com/boginskiy/Clicki/internal/repository"
	valid "github.com/boginskiy/Clicki/internal/validation"
	"github.com/boginskiy/Clicki/pkg"
)

type APIShortURL struct {
	Repo        repo.Repository
	ExtraFuncer prep.ExtraFuncer
	Kwargs      conf.VarGetter
	Checker     valid.Checker
	Logger      logg.Logger
}

func NewAPIShortURL(
	kwargs conf.VarGetter, logger logg.Logger, repo repo.Repository,
	checker valid.Checker, extraFuncer prep.ExtraFuncer) *APIShortURL {

	return &APIShortURL{
		ExtraFuncer: extraFuncer,
		Checker:     checker,
		Logger:      logger,
		Kwargs:      kwargs,
		Repo:        repo,
	}
}

func (s *APIShortURL) encryptionLongURL() (correlID string) {
	for {
		correlID = pkg.Scramble(LONG)                   // Вызов шифратора
		if s.Repo.CheckUnic(context.TODO(), correlID) { // Проверка на уникальность
			break
		}
	}
	return correlID
}

func (s *APIShortURL) GetHeader() string {
	return "application/json"
}

func (s *APIShortURL) Create(req *http.Request) ([]byte, error) {
	// Deserialization Body
	bodyJSON := mod.NewURLJson()
	err := s.ExtraFuncer.Deserialization(req, bodyJSON)

	if err != nil {
		s.Logger.RaiseFatal(err, DeserializFatal, nil)
		return EmptyByteSlice, err
	}

	// Валидируем URL. Проверка регуляркой, что строка является доменом сайта
	if !s.Checker.CheckUpURL(bodyJSON.URL) || bodyJSON.URL == "" {
		s.Logger.RaiseInfo("APIShortURL.Create>CheckUpURL",
			logg.Fields{"error": ErrDataNotValid.Error()})
		return EmptyByteSlice, ErrDataNotValid
	}

	correlationID := s.encryptionLongURL()                              // Уникальный идентификатор
	shortURL := s.Kwargs.GetBaseURL() + "/" + correlationID             // Создаем новый сокращенный URL
	preRecord := mod.NewURLTb(0, correlationID, bodyJSON.URL, shortURL) // Создаем черновую запись
	record, err := s.Repo.Create(context.TODO(), preRecord)             // Кладем в DB данные

	if err != nil && record == nil {
		s.Logger.RaiseError(err, "APIShortURL.Create>Repo.Create", nil)
		return EmptyByteSlice, err
	}

	//
	var resJSON *mod.ResultJSON
	switch r := record.(type) {
	case *mod.URLTb:
		resJSON = mod.NewResultJSON(bodyJSON, r.ShortURL)
	case string:
		resJSON = mod.NewResultJSON(bodyJSON, r)
	default:
		s.Logger.RaiseError(err, "APIShortURL.Create>switch", nil)
		return EmptyByteSlice, err
	}

	// Serialization Body
	result, err2 := s.ExtraFuncer.Serialization(resJSON)

	if err2 != nil {
		s.Logger.RaiseError(err2, "APIShortURL.Create>NewResultJSON", nil)
		return EmptyByteSlice, err2
	}
	return result, err
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
			logg.Fields{"fatal": ErrDataNotValid.Error()})
		return EmptyByteSlice, ErrDataNotValid
	}

	// Разбор тела запроса
	respURLSet := make([]mod.ResURLSet, 0, 10)

	for decoder.More() {
		var rURL mod.ReqURLSet
		err := decoder.Decode(&rURL)

		if err != nil {
			s.Logger.RaiseFatal(err, "ShortURL>SetBatch>Decode", nil)
			return EmptyByteSlice, err
		}

		// TODO!
		// shortURL := s.Kwargs.GetBaseURL() + "/" + s.encryptionLongURL()
		shortURL := s.Kwargs.GetBaseURL() + "/" + rURL.CorrelationID

		// Сбор множества URL
		respURLSet = append(respURLSet, mod.NewResURLSet(rURL.CorrelationID, rURL.OriginalURL, shortURL))
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
