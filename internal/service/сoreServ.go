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

func (c *CoreService) TakeUserIDFromCtx(req *http.Request) int {
	UserID, ok := req.Context().Value(auth.CtxUserID).(int)
	if !ok || UserID <= 0 {
		c.Logg.RaiseError(ErrUserIDNotValid, "CoreService.TakeUserIDFromCtx>CtxUserID", nil)
	}
	return UserID
}

func (c *CoreService) EncrypOriginURL() (correlID string) {
	for {
		correlID = pkg.Scramble(LONG)                         // Вызов шифратора
		if c.Repo.CheckUnicRecord(context.TODO(), correlID) { // Проверка на уникальность
			break
		}
	}
	return correlID
}

func (c *CoreService) EventOfAudit(action string, userID int, url string) {
	// Собираем событие аудита
	event := audit.NewEvent(action, userID, url)
	// Отправка события подписчикам
	c.Publisher.Send(event)
}
