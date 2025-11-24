package auth

import (
	"context"
	"errors"
	"net/http"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
	repo "github.com/boginskiy/Clicki/internal/repository"
)

type Auth struct {
	Cfg        conf.Config
	Logg       logg.Logger
	Repo       repo.Repository
	LastUser   int
	JWTService JWTer
}

func NewAuth(config conf.Config, logger logg.Logger, repo repo.Repository) *Auth {
	return &Auth{
		Cfg:        config,
		Logg:       logger,
		Repo:       repo,
		JWTService: NewJWTService(config),

		// Louding last userID from last record
		LastUser: repo.ReadLastRecord(context.TODO()),
	}
}

func (a *Auth) createCookie(token, name string) *http.Cookie {
	cokiTime := a.Cfg.GetCokiLiveTime()
	return &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		HttpOnly: true,                    // Доступ только серверу, увеличивает безопасность
		SameSite: http.SameSiteStrictMode, // Запрещает отправлять куки с другого домена
		MaxAge:   cokiTime,                // Срок жизни куки
		Secure:   false,                   // Поставьте true, если работаете через HTTPS
	}
}

func (a *Auth) nextUser() int {
	a.LastUser += 1
	return a.LastUser
}

func (a *Auth) Authorization(req *http.Request) (*http.Cookie, int, error) {
	UserID := a.nextUser()
	token, err := a.JWTService.CreateJWT(UserID)
	if err != nil {
		a.Logg.RaiseError(err, "Auth>Authorization>CreateJWT", nil)
		return nil, 0, ErrCreateToken
	}
	return a.createCookie(token, a.Cfg.GetNameCoki()), UserID, nil
}

func (a *Auth) Authentication(req *http.Request) (*http.Cookie, int, error) {
	// Достаем 'Cookie'
	cookie, err := req.Cookie(a.Cfg.GetNameCoki())

	// Авторизация если отсутствуют 'Cookie'
	if err != nil {
		return a.Authorization(req)
	}

	// Аутентификация если присутствуют 'Cookie'
	UserID, err := a.JWTService.GetIDAndValidJWT(cookie.Value)

	// Условие непрохождения аутентификации. Пользователь не найден
	if UserID <= 0 {
		a.Logg.RaiseInfo(ErrUserNotFound.Error(), logg.Fields{"userID": UserID})
		return nil, 0, ErrUserNotFound
	}

	// Условие для обновления токена
	if err != nil {

		if errors.Is(err, ErrTokenIsExpired) || errors.Is(err, ErrTokenNotValid) {
			token, err := a.JWTService.CreateJWT(UserID)
			if err != nil {
				a.Logg.RaiseError(err, "Auth>Authentication>CreateJWT", nil)
				return nil, 0, ErrCreateToken
			}
			// Выдаем свежий токен
			return a.createCookie(token, a.Cfg.GetNameCoki()), UserID, nil
		}
		a.Logg.RaiseError(err, "Auth>Authentication>CreateJWT", nil)
		return nil, 0, ErrTokenNotValid
	}
	return cookie, UserID, nil
}
