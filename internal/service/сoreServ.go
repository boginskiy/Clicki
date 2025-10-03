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
	Logg   logg.Logger
}

func NewCoreService(kwargs conf.VarGetter, logg logg.Logger, repo repository.Repository) *CoreService {
	return &CoreService{
		Logg:   logg,
		Kwargs: kwargs,
		Repo:   repo,
	}
}

func (c *CoreService) TakeUserIDFromCtx(req *http.Request) int {
	UserID, ok := req.Context().Value(midw.CtxUserID).(int)
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
