package service

import (
	"context"
	"net/http"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
	midw "github.com/boginskiy/Clicki/internal/middleware"
	"github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/pkg"
)

type CoreService struct {
	Repo   repository.Repository
	Kwargs conf.VarGetter
	Logger logg.Logger
}

func NewCoreService(kwargs conf.VarGetter, logger logg.Logger, repo repository.Repository) CoreServicer {
	return &CoreService{
		Logger: logger,
		Kwargs: kwargs,
		Repo:   repo,
	}
}

func (c *CoreService) takeUserIdFromCtx(req *http.Request) int {
	UserID, ok := req.Context().Value(midw.CtxUserID).(int)
	if !ok || UserID <= 0 {
		c.Logger.RaiseError(ErrUserIDNotValid, "CoreService.takeUserIdFromCtx>CtxUserID", nil)
		UserID = -1
	}
	return UserID
}

func (s *CoreService) encrypOriginURL() (correlID string) {
	for {
		correlID = pkg.Scramble(LONG)                   // Вызов шифратора
		if s.Repo.CheckUnic(context.TODO(), correlID) { // Проверка на уникальность
			break
		}
	}
	return correlID
}
