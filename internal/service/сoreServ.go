package service

import (
	"context"
	"net/http"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/audit"
	"github.com/boginskiy/Clicki/internal/auth"
	"github.com/boginskiy/Clicki/internal/logg"
	repo "github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/pkg"
)

// CoreService - core.
type CoreService struct {
	Repo      repo.Repository
	Publisher audit.Publisher
	Cfg       conf.Config
	Logg      logg.Logger
}

func NewCoreService(
	config conf.Config,
	logger logg.Logger,
	repository repo.Repository,
	publisher audit.Publisher) *CoreService {

	return &CoreService{
		Logg:      logger,
		Cfg:       config,
		Repo:      repository,
		Publisher: publisher,
	}
}

// TakeUserIDFromCtx -
func (c *CoreService) TakeUserIDFromCtx(req *http.Request) int {
	UserID, ok := req.Context().Value(auth.CtxUserID).(int)
	if !ok || UserID <= 0 {
		c.Logg.RaiseError(ErrUserIDNotValid, "CoreService.TakeUserIDFromCtx>CtxUserID", nil)
	}
	return UserID
}

// EncrypOriginURL -
func (c *CoreService) EncrypOriginURL() (correlID string) {
	for {
		correlID = pkg.Scramble(LONG)                         // Call scramble.
		if c.Repo.CheckUnicRecord(context.TODO(), correlID) { // Check on unic.
			break
		}
	}
	return correlID
}

// EventOfAudit -
func (c *CoreService) EventOfAudit(action string, userID int, url string) {
	// Collection event audit.
	event := audit.NewEvent(action, userID, url)
	// Send event.
	if c.Publisher != nil {
		c.Publisher.Send(event)
	}
}
