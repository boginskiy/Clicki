package service

import (
	"context"
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

func (s *URLServ) GetStats(req *http.Request) ([]byte, error) {
	return StoreDBIsSucces, nil
}

func (s *URLServ) ReadSet(ctx context.Context, protocol p.Protocol) (any, error) {
	return EmptyByteSlice, nil
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

	modURLTb := model.NewURLTb(0, correlationID, originURL, URLServ, userID) // Create record.
	record, err := s.Repo.CreateRecord(context.TODO(), modURLTb)             // Put record in the DB.

	if record == nil {
		s.Logg.RaiseError(err, "URLServ.CreateURL>Repo.Create", nil)
		return EmptyByteSlice, err
	}

	// Audit.
	s.eventOfAudit("shorten", userID, originURL)

	return []byte(record.ShortURL), err
}

func (s *URLServ) Read(ctx context.Context, protocol p.Protocol, request any) ([]byte, error) {
	// Take user id.
	userID, err := s.getUserIDFromCtx(ctx)
	if err != nil {
		return EmptyByteSlice, err
	}

	// Take params correlationID.
	correlationID, err := protocol.GetURLID(request)
	if err != nil {
		return EmptyByteSlice, err
	}

	// Take origin URL.
	record, err := s.Repo.ReadRecord(context.TODO(), correlationID)
	if err != nil {
		s.Logg.RaiseError(err, "URLServ.Read>DB.Read", nil)
		return EmptyByteSlice, ErrDataNotValid
	}

	// if flag == true, record is in queue on deleting
	if record.DeletedFlag {
		return EmptyByteSlice, ErrReadRecord
	}

	// Audit.
	s.eventOfAudit("follow", userID, record.OriginalURL)

	return []byte(record.OriginalURL), nil
}

func (s *URLServ) takeUserIDFromCtx(req *http.Request) int {
	UserID, ok := req.Context().Value(auth.CtxUserID).(int)
	if !ok || UserID <= 0 {
		s.Logg.RaiseError(ErrUserIDNotValid, "URLServ.takeUserIDFromCtx>CtxUserID", nil)
	}
	return UserID
}

func (s *URLServ) getUserIDFromCtx(ctx context.Context) (int, error) {
	var userID int
	UserID, ok := ctx.Value(auth.CtxUserID).(int)
	if !ok || UserID <= 0 {
		return userID, ErrUserIDNotValid
	}
	return UserID, nil
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
