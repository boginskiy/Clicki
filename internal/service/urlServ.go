package service

import (
	"context"
	"net/http"
	"strings"

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

// URLServ - service about generation short URL.
type URLServ struct {
	Cfg       conf.Config
	Logg      logg.Logger
	Funcer    prep.Funcer
	Repo      repo.Repository
	Checker   validation.Checker
	Publisher audit.Publisher
}

func NewURLServ(
	config conf.Config,
	logger logg.Logger,
	repository repo.Repository,
	checker validation.Checker,
	fancer prep.Funcer,
	publisher audit.Publisher) *URLServ {

	return &URLServ{
		Cfg:       config,
		Logg:      logger,
		Funcer:    fancer,
		Checker:   checker,
		Repo:      repository,
		Publisher: publisher,
	}
}

func (s *URLServ) CreateSet(req *http.Request) ([]byte, error) {
	return StoreDBIsSucces, nil
}

func (s *URLServ) ReadSet(req *http.Request) ([]byte, error) {
	return StoreDBIsSucces, nil
}

func (s *URLServ) GetStats(req *http.Request) ([]byte, error) {
	return StoreDBIsSucces, nil
}

func (s *URLServ) CheckDB(req *http.Request) ([]byte, error) {
	_, err := s.Repo.PingStore(context.TODO())
	if err != nil {
		s.Logg.RaiseFatal(err, "URLServ.CreaCheckPingte>Ping", nil)
		return EmptyByteSlice, err
	}
	return StoreDBIsSucces, nil
}

func (s *URLServ) Create(req *http.Request) ([]byte, error) {
	// Take body request.
	originURL, err := s.Funcer.TakeAllBodyFromReq(req)
	if err != nil {
		s.Logg.RaiseFatal(err, "URLServ.CreateURL>TakeAllBodyFromReq", nil)
		return EmptyByteSlice, err
	}

	// Validation URL. Check regular expression, that line is domen of site.
	if !s.Checker.CheckUpURL(originURL) || originURL == "" {
		s.Logg.RaiseError(ErrDataNotValid, "URLServ.CreateURL>CheckUpURL", nil)
		return EmptyByteSlice, ErrDataNotValid
	}

	userID := s.takeUserIDFromCtx(req)                  // Take user id.
	correlationID := s.encrypOriginURL()                // Take unic id.
	URLServ := s.Cfg.GetBaseURL() + "/" + correlationID // New short URL.

	preRecord := model.NewURLTb(0, correlationID, originURL, URLServ, userID) // Create record.
	record, err := s.Repo.CreateRecord(context.TODO(), preRecord)             // Put record in the DB.

	if err != nil && record == nil {
		s.Logg.RaiseError(err, "URLServ.CreateURL>Repo.Create", nil)
		return EmptyByteSlice, err
	}

	// Audit.
	s.eventOfAudit("shorten", userID, originURL)

	return []byte(record.(*model.URLTb).ShortURL), err
}

func (s *URLServ) Read(req *http.Request) ([]byte, error) {
	userID := s.takeUserIDFromCtx(req)                              // Take user id.
	correlationID := strings.TrimLeft(req.URL.Path, "/")            // Take params correlationID.
	record, err := s.Repo.ReadRecord(context.TODO(), correlationID) // Take origin URL.

	if err != nil {
		s.Logg.RaiseError(err, "URLServ.Read>DB.Read", nil)
		return EmptyByteSlice, ErrDataNotValid
	}

	if r, ok := record.(*model.URLTb); ok {
		// if flag == true, record is in queue on deleting
		if r.DeletedFlag {
			return EmptyByteSlice, ErrReadRecord

		} else {
			// Audit.
			s.eventOfAudit("follow", userID, r.OriginalURL)
			return []byte(r.OriginalURL), nil
		}
	}

	// Default.
	s.Logg.RaiseError(err, "URLServ.Read>DB.Read>switch", nil)
	return EmptyByteSlice, ErrDataNotValid
}

func (s *URLServ) takeUserIDFromCtx(req *http.Request) int {
	UserID, ok := req.Context().Value(auth.CtxUserID).(int)
	if !ok || UserID <= 0 {
		s.Logg.RaiseError(ErrUserIDNotValid, "URLServ.takeUserIDFromCtx>CtxUserID", nil)
	}
	return UserID
}

func (s *URLServ) encrypOriginURL() (correlID string) {
	for {
		correlID = pkg.Scramble(LONG)                           // Call scramble.
		if s.Repo.CheckUniqueRecord(context.TODO(), correlID) { // Check on unic.
			break
		}
	}
	return correlID
}

func (s *URLServ) eventOfAudit(action string, userID int, url string) {
	// Collection event audit.
	event := audit.NewEvent(action, userID, url)
	// Send event.
	if s.Publisher != nil {
		s.Publisher.Send(event)
	}
}
