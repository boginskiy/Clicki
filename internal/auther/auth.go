package auther

import (
	"context"
	"net/http"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
	repo "github.com/boginskiy/Clicki/internal/repository"
)

type Auth struct {
	Kwargs   conf.VarGetter
	Logger   logg.Logger
	Repo     repo.Repository
	LastUser int
	JWTService
}

func NewAuth(kwargs conf.VarGetter, logger logg.Logger, repo repo.Repository) *Auth {
	return &Auth{
		Kwargs: kwargs,
		Logger: logger,
		Repo:   repo,

		// Подгружаем данные из Репозитория об ID последнего записанного пользователя
		// По ходу работы программы, новые пользователи будут получать следующие по порядку ID
		LastUser: repo.ReadLastRecord(context.TODO()),
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
	a.LastUser += 1
	return a.LastUser
}
