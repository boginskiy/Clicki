package service

import (
	"context"
	"net/http"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/audit"
	"github.com/boginskiy/Clicki/internal/auther"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/pkg"
)

type CoreService struct {
	Repo      repository.Repository
	Publisher audit.Publisher
	Kwargs    conf.VarGetter
	Logg      logg.Logger
}

func NewCoreService(
	kwargs conf.VarGetter,
	logg logg.Logger,
	repo repository.Repository,
	publisher audit.Publisher) *CoreService {

	return &CoreService{
		Logg:      logg,
		Kwargs:    kwargs,
		Repo:      repo,
		Publisher: publisher,
	}
}

func (c *CoreService) TakeUserIDFromCtx(req *http.Request) int {
	UserID, ok := req.Context().Value(auther.CtxUserID).(int)
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
