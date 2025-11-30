package service

import (
	"context"
	"encoding/json"
	"net/http"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/audit"
	"github.com/boginskiy/Clicki/internal/auth"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/model"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	repo "github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/internal/validation"
	"github.com/boginskiy/Clicki/pkg"
)

// APIURLServ - struct about service api.
type APIURLServ struct {
	Cfg       conf.Config
	Logg      logg.Logger
	Repo      repo.Repository
	Funcer    prep.Funcer
	Checker   validation.Checker
	Publisher audit.Publisher
}

func NewAPIURLServ(
	config conf.Config,
	logger logg.Logger,
	repository repo.Repository,
	checker validation.Checker,
	fancer prep.Funcer,
	publisher audit.Publisher) *APIURLServ {

	return &APIURLServ{
		Cfg:       config,
		Logg:      logger,
		Funcer:    fancer,
		Checker:   checker,
		Repo:      repository,
		Publisher: publisher,
	}
}

func (s *APIURLServ) Read(req *http.Request) ([]byte, error) {
	return EmptyByteSlice, nil
}

func (s *APIURLServ) CheckDB(req *http.Request) ([]byte, error) {
	return EmptyByteSlice, nil
}

func (s *APIURLServ) Create(req *http.Request) ([]byte, error) {
	// Deserialization Body.
	bodyJSON := model.NewURLJson()
	err := s.Funcer.Deserialization(req, bodyJSON)

	if err != nil {
		s.Logg.RaiseFatal(err, "deserialization was bad", nil)
		return EmptyByteSlice, err
	}

	// Валидируем URL. Проверка регуляркой, что строка является доменом сайта.
	if !s.Checker.CheckUpURL(bodyJSON.URL) || bodyJSON.URL == "" {
		s.Logg.RaiseInfo("APIURLServ.Create>CheckUpURL",
			logg.Fields{"error": ErrDataNotValid.Error()})
		return EmptyByteSlice, ErrDataNotValid
	}

	userID := s.takeUserIDFromCtx(req)                   // Take id user.
	correlationID := s.encrypOriginURL()                 // Create unic id.
	shortURL := s.Cfg.GetBaseURL() + "/" + correlationID // Create new short URL.

	preRecord := model.NewURLTb(0, correlationID, bodyJSON.URL, shortURL, userID) // Create dirty record.
	record, err := s.Repo.CreateRecord(context.TODO(), preRecord)                 // Put record in the DB.

	if err != nil && record == nil {
		s.Logg.RaiseError(err, "APIURLServ.Create>Repo.Create", nil)
		return EmptyByteSlice, err
	}

	// Audition.
	s.eventOfAudit("shorten", userID, bodyJSON.URL)

	// Definition type.
	var resJSON *model.ResultJSON
	switch r := record.(type) {
	case *model.URLTb:
		resJSON = model.NewResultJSON(bodyJSON, r.ShortURL)
	case string:
		resJSON = model.NewResultJSON(bodyJSON, r)
	default:
		s.Logg.RaiseError(err, "APIURLServ.Create>switch", nil)
		return EmptyByteSlice, err
	}

	// Serialization.
	result, err2 := s.Funcer.Serialization(resJSON)
	if err2 != nil {
		s.Logg.RaiseError(err2, "APIURLServ.Create>NewResultJSON", nil)
		return EmptyByteSlice, err2
	}
	return result, err
}

func (s *APIURLServ) CreateSet(req *http.Request) ([]byte, error) {
	// Create decoder.
	decoder := json.NewDecoder(req.Body)

	// Check, that right thing has arrived.
	token, _ := decoder.Token()
	if token != json.Delim('[') {
		s.Logg.RaiseFatal(ErrDataNotValid, "ShortURL>SetBatch>Token",
			logg.Fields{"fatal": ErrDataNotValid.Error()})
		return EmptyByteSlice, ErrDataNotValid
	}

	// Take id user.
	userID := s.takeUserIDFromCtx(req)

	// for parsing body of request.
	respURLSet := make([]model.ResURLSet, 0, 10)

	for decoder.More() {
		var rURL model.ReqURLSet
		err := decoder.Decode(&rURL)

		if err != nil {
			s.Logg.RaiseFatal(err, "ShortURL>SetBatch>Decode", nil)
			return EmptyByteSlice, err
		}

		shortURL := s.Cfg.GetBaseURL() + "/" + rURL.CorrelationID
		// Collection set URL.
		respURLSet = append(respURLSet, model.NewResURLSet(
			rURL.CorrelationID, rURL.OriginalURL, shortURL, userID))
	}

	// Save in the DB.
	err := s.Repo.CreateRecords(context.TODO(), respURLSet)
	if err != nil {
		s.Logg.RaiseError(err, "APIURLServ>SetBatch>CreateSet", nil)
		return EmptyByteSlice, err
	}

	// Serialization.
	result, err := json.Marshal(respURLSet)
	s.Logg.RaiseFatal(err, "ShortURL>SetBatch>Marshal", nil)
	return result, nil
}

func (s *APIURLServ) ReadSet(req *http.Request) ([]byte, error) {
	// Take id user.
	userID := s.takeUserIDFromCtx(req)

	dataSet, err := s.Repo.ReadRecords(context.TODO(), userID)
	if err != nil {
		s.Logg.RaiseError(err, "APIURLServ.ReadSetUserURL>ReadSet", nil)
		return EmptyByteSlice, err
	}

	// Definition record by user.
	records, ok := dataSet.([]model.ResUserURLSet)
	if !ok {
		s.Logg.RaiseError(ErrDataNotValid, "APIURLServ.ReadSetUserURL>Type?", nil)
		return EmptyByteSlice, ErrDataNotValid
	}
	if len(records) == 0 {
		return EmptyByteSlice, nil
	}

	// Serialization.
	result, err := s.Funcer.Serialization(records)
	if err != nil {
		s.Logg.RaiseError(err, "APIURLServ.ReadSetUserURL>Serialization", nil)
		return EmptyByteSlice, err
	}
	return result, err
}

func (s *APIURLServ) takeUserIDFromCtx(req *http.Request) int {
	UserID, ok := req.Context().Value(auth.CtxUserID).(int)
	if !ok || UserID <= 0 {
		s.Logg.RaiseError(ErrUserIDNotValid, "APIURLServ.takeUserIDFromCtx>CtxUserID", nil)
	}
	return UserID
}

func (s *APIURLServ) encrypOriginURL() (correlID string) {
	for {
		correlID = pkg.Scramble(LONG)                           // Call scramble.
		if s.Repo.CheckUniqueRecord(context.TODO(), correlID) { // Check on unic.
			break
		}
	}
	return correlID
}

func (s *APIURLServ) eventOfAudit(action string, userID int, url string) {
	// Collection event audit.
	event := audit.NewEvent(action, userID, url)
	// Send event.
	if s.Publisher != nil {
		s.Publisher.Send(event)
	}
}
