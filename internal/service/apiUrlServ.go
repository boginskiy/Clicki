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
	p "github.com/boginskiy/Clicki/internal/protocol"
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

func (s *APIURLServ) GetStats(req *http.Request) ([]byte, error) {
	// In this place we know about who need stats
	tmpMap := map[string]int{
		"urls":  s.Repo.ReadQuantityShortURLs(context.TODO()), // quantityURLs
		"users": s.Repo.ReadQuantityUsers(context.TODO()),     // quantityUsers
	}
	return s.Funcer.Serialization(tmpMap), nil
}

func (s *APIURLServ) Read(req *http.Request) ([]byte, error) {
	return EmptyByteSlice, nil
}

func (s *APIURLServ) CheckDB(req *http.Request) ([]byte, error) {
	return EmptyByteSlice, nil
}

func (s *APIURLServ) Create(ctx context.Context, protocol p.Protocol, request any) ([]byte, error) {
	// Take URL from request.
	urlJSON, err := protocol.GetURLFromRequest(request)
	if err != nil {
		return EmptyByteSlice, err
	}

	// Validation URL.
	if !s.Checker.CheckUpURL(urlJSON.URL) || urlJSON.URL == "" {
		s.Logg.RaiseError(ErrDataNotValid, "CheckUpURL", nil)
		return EmptyByteSlice, ErrDataNotValid
	}

	// Take userID from context.
	userID, err := protocol.GetUserIDFromCtx(ctx)
	if err != nil {
		return EmptyByteSlice, err
	}

	correlationID := s.encrypOriginURL()                 // Create unic id.
	shortURL := s.Cfg.GetBaseURL() + "/" + correlationID // Create new short URL.

	modURLTb := model.NewURLTb(0, correlationID, urlJSON.URL, shortURL, userID) // Create dirty record.
	record, errDB := s.Repo.CreateRecord(context.TODO(), modURLTb)              // Put record in the DB.

	if record == nil {
		s.Logg.RaiseError(errDB, "APIURLServ.Create>Repo.Create", nil)
		return EmptyByteSlice, errDB
	}

	// Audition.
	s.eventOfAudit("shorten", userID, urlJSON.URL)

	// Result.
	resJSON := model.NewResultJSON(urlJSON, record.ShortURL)

	return s.Funcer.Serialization(resJSON), errDB
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

	return s.Funcer.Serialization(respURLSet), nil
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

	return s.Funcer.Serialization(records), nil
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
