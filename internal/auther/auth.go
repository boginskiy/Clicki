package auther

import (
	"context"
	"net/http"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
	repo "github.com/boginskiy/Clicki/internal/repository"
)

type Auth struct {
	Kwargs conf.VarGetter
	Logger logg.Logger
	Repo   repo.Repository
	JWTService
}

func NewAuth(kwargs conf.VarGetter, logger logg.Logger, repo repo.Repository) *Auth {
	return &Auth{
		Kwargs: kwargs,
		Logger: logger,
		Repo:   repo,
	}
}

func (a *Auth) CreateCookie(token, name string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		HttpOnly: true,                    // Доступ только серверу, увеличивает безопасность
		SameSite: http.SameSiteStrictMode, // Запрещает отправлять куки с другого домена
		MaxAge:   300,                     // Жива 300 секунд
		Secure:   false,                   // Поставьте true, если работаете через HTTPS
	}
}

func (a *Auth) NextUser() int {
	lastUser, err := a.Repo.TakeLastUser(context.TODO())
	if err != nil {
		a.Logger.RaiseError(err, "Auth>NextUser>TakeLastUser", nil)
		return lastUser // default value == 1
	}
	return (lastUser + 1)
}

func (a *Auth) CheckUser(userID int) (bool, error) {
	return a.Repo.CheckUser(context.TODO(), userID)
}
