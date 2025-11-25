package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/model"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	repo "github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/internal/validation"
)

// APIShortURL - struct about service api.
type APIShortURL struct {
	Repo    repo.Repository
	ExFunc  prep.ExtraFuncer
	Checker validation.Checker
	Core    *CoreService
}

func NewAPIShortURL(
	core *CoreService,
	repository repo.Repository,
	checker validation.Checker,
	extraFuncer prep.ExtraFuncer) *APIShortURL {

	return &APIShortURL{
		ExFunc:  extraFuncer,
		Checker: checker,
		Core:    core,
		Repo:    repository,
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
	// Deserialization Body.
	bodyJSON := model.NewURLJson()
	err := s.ExFunc.Deserialization(req, bodyJSON)

	if err != nil {
		s.Core.Logg.RaiseFatal(err, DeserializFatal, nil)
		return EmptyByteSlice, err
	}

	// Валидируем URL. Проверка регуляркой, что строка является доменом сайта.
	if !s.Checker.CheckUpURL(bodyJSON.URL) || bodyJSON.URL == "" {
		s.Core.Logg.RaiseInfo("APIShortURL.Create>CheckUpURL",
			logg.Fields{"error": ErrDataNotValid.Error()})
		return EmptyByteSlice, ErrDataNotValid
	}

	userID := s.Core.TakeUserIDFromCtx(req)                   // Take id user.
	correlationID := s.Core.EncrypOriginURL()                 // Create unic id.
	shortURL := s.Core.Cfg.GetBaseURL() + "/" + correlationID // Create new short URL.

	preRecord := model.NewURLTb(0, correlationID, bodyJSON.URL, shortURL, userID) // Create dirty record.
	record, err := s.Repo.CreateRecord(context.TODO(), preRecord)                 // Put record in the DB.

	if err != nil && record == nil {
		s.Core.Logg.RaiseError(err, "APIShortURL.Create>Repo.Create", nil)
		return EmptyByteSlice, err
	}

	// Audition.
	s.Core.EventOfAudit("shorten", userID, bodyJSON.URL)

	// Definition type.
	var resJSON *model.ResultJSON
	switch r := record.(type) {
	case *model.URLTb:
		resJSON = model.NewResultJSON(bodyJSON, r.ShortURL)
	case string:
		resJSON = model.NewResultJSON(bodyJSON, r)
	default:
		s.Core.Logg.RaiseError(err, "APIShortURL.Create>switch", nil)
		return EmptyByteSlice, err
	}

	// Serialization.
	result, err2 := s.ExFunc.Serialization(resJSON)
	if err2 != nil {
		s.Core.Logg.RaiseError(err2, "APIShortURL.Create>NewResultJSON", nil)
		return EmptyByteSlice, err2
	}
	return result, err
}

func (s *APIShortURL) CreateSetURL(req *http.Request) ([]byte, error) {
	// Create decoder.
	decoder := json.NewDecoder(req.Body)

	// Check, that right thing has arrived.
	token, _ := decoder.Token()
	if token != json.Delim('[') {
		s.Core.Logg.RaiseFatal(ErrDataNotValid, "ShortURL>SetBatch>Token",
			logg.Fields{"fatal": ErrDataNotValid.Error()})
		return EmptyByteSlice, ErrDataNotValid
	}

	// Take id user.
	userID := s.Core.TakeUserIDFromCtx(req)

	// for parsing body of request.
	respURLSet := make([]model.ResURLSet, 0, 10)

	for decoder.More() {
		var rURL model.ReqURLSet
		err := decoder.Decode(&rURL)

		if err != nil {
			s.Core.Logg.RaiseFatal(err, "ShortURL>SetBatch>Decode", nil)
			return EmptyByteSlice, err
		}

		shortURL := s.Core.Cfg.GetBaseURL() + "/" + rURL.CorrelationID
		// Collection set URL.
		respURLSet = append(respURLSet, model.NewResURLSet(
			rURL.CorrelationID, rURL.OriginalURL, shortURL, userID))
	}

	// Save in the DB.
	err := s.Repo.CreateRecords(context.TODO(), respURLSet)
	if err != nil {
		s.Core.Logg.RaiseError(err, "APIShortURL>SetBatch>CreateSet", nil)
		return EmptyByteSlice, err
	}

	// Serialization.
	result, err := json.Marshal(respURLSet)
	s.Core.Logg.RaiseFatal(err, "ShortURL>SetBatch>Marshal", nil)
	return result, nil
}

func (s *APIShortURL) ReadSetUserURL(req *http.Request) ([]byte, error) {
	// Take id user.
	userID := s.Core.TakeUserIDFromCtx(req)

	dataSet, err := s.Repo.ReadRecords(context.TODO(), userID)
	if err != nil {
		s.Core.Logg.RaiseError(err, "APIShortURL.ReadSetUserURL>ReadSet", nil)
		return EmptyByteSlice, err
	}

	// Definition record by user.
	records, ok := dataSet.([]model.ResUserURLSet)
	if !ok {
		s.Core.Logg.RaiseError(ErrDataNotValid, "APIShortURL.ReadSetUserURL>Type?", nil)
		return EmptyByteSlice, ErrDataNotValid
	}
	if len(records) == 0 {
		return EmptyByteSlice, nil
	}

	// Serialization.
	result, err := s.ExFunc.Serialization(records)
	if err != nil {
		s.Core.Logg.RaiseError(err, "APIShortURL.ReadSetUserURL>Serialization", nil)
		return EmptyByteSlice, err
	}
	return result, err
}
