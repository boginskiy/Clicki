package register

import (
	"context"
	"errors"
	"net/http"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/model"
	"github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/internal/user"
)

/*
Steps:
	Зарегистрироваться → Идентифицироваться → Aутентифицироваться → Авторизироваться
*/

type Regist struct {
	Kwargs   conf.VarGetter
	Logger   logg.Logger
	RepoUser repository.RepositoryUser
	JWTService
}

func NewRegist(kwargs conf.VarGetter, logger logg.Logger, repoUser repository.RepositoryUser) *Regist {
	// Аутентификация пользователя

	return &Regist{
		Kwargs:   kwargs,
		Logger:   logger,
		RepoUser: repoUser,
	}
}

func (r *Regist) CreateCookie(token, name string) *http.Cookie {
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

func (r *Regist) createNewUser() (int, error) {
	// Пока считаем такого пользователя новым
	newUser := user.NewUser().CreateEmpty()
	// Записываем в БД пользователя, получаем userID
	return r.RepoUser.CreateUser(context.TODO(), newUser)
}

func (r *Regist) checkUser(userID int) (bool, error) {
	record := r.RepoUser.ReadUser(context.TODO(), userID)
	user, ok := record.(*model.UserTb)
	if !ok {
		return false, ErrDataNotValid
	}
	if 0 < user.ID {
		return true, nil
	}
	return false, nil
}

func (r *Regist) Registration(req *http.Request, userID *int) ([]byte, *http.Cookie, error) {
	// Регистрация пользователя посредством Куки
	_, err := req.Cookie(NameCoki)

	// У пользователя отсутствует 'Cookie'
	if err != nil {
		// Создаем нового пользователя
		*userID, err = r.createNewUser()

		if err != nil {
			r.Logger.RaiseError(err, "Regist>Authentication>createNewUser", nil)
			return MessRegistIsBad, nil, err
		}
		// Создаем новый токен для нового пользователя
		token, err := r.CreateJWT(*userID)
		if err != nil {
			r.Logger.RaiseError(err, "Regist>Authentication>CreateJWT", nil)
			return MessRegistIsBad, nil, err
		}
		// Отправляем куки пользователю
		return MessRegistCompleted, r.CreateCookie(token, NameCoki), nil
	}
	return nil, nil, nil
}

func (r *Regist) Authentication(req *http.Request, userID *int) ([]byte, *http.Cookie, error) {
	// Аутентификация пользователя через Куки
	cookie, _ := req.Cookie(NameCoki)

	// У пользователя присутствует 'Cookie'
	var err error
	*userID, err = r.GetIDAndValidJWT(cookie.Value)

	// Определеяем условие для обновления токена
	updateToken := errors.Is(err, ErrTokenIsExpired) || errors.Is(err, ErrTokenNotValid)

	if err != nil && !updateToken {
		r.Logger.RaiseError(err, "Regist>Authentication.GetIDAndValidJWT", nil)
		return MessProcessingIsBad, nil, err
	}

	// Если пользователь есть в БД и его токен не валидный или просрочен
	// Обновляем ему токен
	isThereUser, err := r.checkUser(*userID)
	if err != nil {
		r.Logger.RaiseError(err, "Regist>Authentication.checkUser", nil)
	}

	if isThereUser && updateToken {
		// Создаем новый токен для существующего пользователя
		token, err := r.CreateJWT(*userID)
		if err != nil {
			r.Logger.RaiseError(err, "Regist>Authentication>CreateJWT", nil)
			return MessProcessingIsBad, nil, err
		}
		// Отправляем новую куку пользователю
		return MessAuthCompleted2, r.CreateCookie(token, NameCoki), nil
	}

	// Проверяем, что пользователя нет в БД
	if !isThereUser {
		return MessNeedRegist, nil, nil
	}

	// У пользователя есть валидный токен и он присутствует в БД
	return nil, nil, nil
}
