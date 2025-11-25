package service

import (
	"context"
	"net/http"
	"strings"

	"github.com/boginskiy/Clicki/internal/model"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	repo "github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/internal/validation"
)

// ShortURL - service about generation short URL.
type ShortURL struct {
	ExFunc  prep.ExtraFuncer
	Repo    repo.Repository
	Checker validation.Checker
	Core    *CoreService
}

func NewShortURL(
	core *CoreService,
	repository repo.Repository,
	checker validation.Checker,
	extraFuncer prep.ExtraFuncer) *ShortURL {

	return &ShortURL{
		ExFunc:  extraFuncer,
		Core:    core,
		Checker: checker,
		Repo:    repository,
	}
}

func (s *ShortURL) CreateSetURL(req *http.Request) ([]byte, error) {
	return StoreDBIsSucces, nil
}

func (s *ShortURL) ReadSetUserURL(req *http.Request) ([]byte, error) {
	return StoreDBIsSucces, nil
}

func (s *ShortURL) GetHeader() string {
	return "text/plain"
}

func (s *ShortURL) CreateURL(req *http.Request) ([]byte, error) {
	// Take body request.
	originURL, err := s.ExFunc.TakeAllBodyFromReq(req)
	if err != nil {
		s.Core.Logg.RaiseFatal(err, "ShortURL.CreateURL>TakeAllBodyFromReq", nil)
		return EmptyByteSlice, err
	}

	// Validation URL. Check regular expression, that line is domen of site.
	if !s.Checker.CheckUpURL(originURL) || originURL == "" {
		s.Core.Logg.RaiseError(ErrDataNotValid, "ShortURL.CreateURL>CheckUpURL", nil)
		return EmptyByteSlice, ErrDataNotValid
	}

	userID := s.Core.TakeUserIDFromCtx(req)                   //  Take user id.
	correlationID := s.Core.EncrypOriginURL()                 // Take unic id.
	shortURL := s.Core.Cfg.GetBaseURL() + "/" + correlationID // New short URL.

	preRecord := model.NewURLTb(0, correlationID, originURL, shortURL, userID) // Create record.
	record, err := s.Repo.CreateRecord(context.TODO(), preRecord)              // Put record in the DB.

	if err != nil && record == nil {
		s.Core.Logg.RaiseError(err, "ShortURL.CreateURL>Repo.Create", nil)
		return EmptyByteSlice, err
	}

	// Audit.
	s.Core.EventOfAudit("shorten", userID, originURL)

	return []byte(record.(*model.URLTb).ShortURL), err
}

func (s *ShortURL) ReadURL(req *http.Request) ([]byte, error) {
	userID := s.Core.TakeUserIDFromCtx(req)                         // Take user id.
	correlationID := strings.TrimLeft(req.URL.Path, "/")            // Take params correlationID.
	record, err := s.Repo.ReadRecord(context.TODO(), correlationID) // Take origin URL.

	if err != nil {
		s.Core.Logg.RaiseError(err, "ShortURL.Read>DB.Read", nil)
		return EmptyByteSlice, ErrDataNotValid
	}

	if r, ok := record.(*model.URLTb); ok {
		// if flag == true, record is in queue on deleting
		if r.DeletedFlag {
			return EmptyByteSlice, ErrReadRecord

		} else {
			// Audit.
			s.Core.EventOfAudit("follow", userID, r.OriginalURL)
			return []byte(r.OriginalURL), nil
		}
	}

	// Default.
	s.Core.Logg.RaiseError(err, "ShortURL.Read>DB.Read>switch", nil)
	return EmptyByteSlice, ErrDataNotValid
}

func (s *ShortURL) CheckDB(req *http.Request) ([]byte, error) {
	_, err := s.Repo.PingDB(context.TODO())
	if err != nil {
		s.Core.Logg.RaiseFatal(err, "ShortURL.CreaCheckPingte>Ping", nil)
		return EmptyByteSlice, err
	}
	return StoreDBIsSucces, nil
}
