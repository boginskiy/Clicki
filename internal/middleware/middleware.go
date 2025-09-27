package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	auth "github.com/boginskiy/Clicki/internal/auther"
	"github.com/boginskiy/Clicki/internal/gzip"
	"github.com/boginskiy/Clicki/internal/logg"
)

type MvFunc func(http.HandlerFunc) http.HandlerFunc

type Middleware struct {
	Auther auth.Auther
	Logger logg.Logger
}

func NewMiddleware(logger logg.Logger, auther auth.Auther) *Middleware {
	return &Middleware{Logger: logger, Auther: auther}
}

func (m *Middleware) Conveyor(next http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range []MvFunc{m.WithAuth, m.WithInfoLogger, m.WithGzip} {
		next = middleware(next)
	}
	return next
}

func (m *Middleware) WithInfoLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		// Extension standart ResponseWriter
		extW := NewExResWriter(w)
		next(extW, r)

		duration := time.Since(start)

		m.Logger.RaiseInfo(
			logg.DataReqResInfo,
			map[string]any{
				"uri":      uri,
				"method":   method,
				"duration": duration,
				"status":   extW.status,
				"size":     extW.size,
			})
	}
}

func (m *Middleware) WithGzip(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpW := w

		// Checking Encoding and Type
		acceptEncoding := r.Header.Get("Accept-Encoding")
		acceptContent := r.Header.Get("Content-Type")

		jsonGzip := strings.Contains(acceptContent, "application/json")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		htmlGzip := strings.Contains(acceptContent, "text/html")

		if supportsGzip && (jsonGzip || htmlGzip) {
			// Оборачиваем http.ResponseWriter новым с gzip
			compW := gzip.NewCompressWriter(w)
			tmpW = compW
			defer compW.Close()
		}

		// Проверка, что клиент отправил сжатые данные
		contentEncoding := r.Header.Get("Content-Encoding")
		sendGzip := strings.Contains(contentEncoding, "gzip")

		if sendGzip {
			decompR, err := gzip.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// Меняем тело запроса на новое
			r.Body = decompR
		}
		// Передача управления
		next(tmpW, r)
	}
}

func (m *Middleware) WithAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie(NAMECOKI) // Достаем 'Cookie'
		var UserID int                    // Идентификатор пользователя

		// У пользователя отсутствует 'Cookie'. Авторизация
		if err != nil {
			UserID = m.Auther.NextUser()
			token, err := m.Auther.CreateJWT(UserID)
			if err != nil {
				m.Logger.RaiseError(err, "Middleware>WithAuth>CreateJWT", nil)
			}
			log.Println("WithAuth>отсутствует 'Cookie'", UserID)
			cookie := m.Auther.CreateCookie(token, NAMECOKI)
			http.SetCookie(w, cookie)

		} else {
			// У пользователя есть 'Cookie'. Аутентификация
			UserID, err = m.Auther.GetIDAndValidJWT(cookie.Value)

			log.Println("WithAuth>Присутствует 'Cookie'", UserID)

			// Условие непрохождения аутентификации
			if UserID <= 0 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Анализ ошибок
			if err != nil {

				// Условие для обновления токена
				if errors.Is(err, auth.ErrTokenIsExpired) || errors.Is(err, auth.ErrTokenNotValid) {
					token, err := m.Auther.CreateJWT(UserID)
					if err != nil {
						m.Logger.RaiseError(err, "Middleware>WithAuth>CreateJWT", nil)
						// TODO! хреновое место
					}
					// Выдаем свежий токен
					cookie := m.Auther.CreateCookie(token, NAMECOKI)
					http.SetCookie(w, cookie)

				} else {
					m.Logger.RaiseError(err, "Middleware>WithAuth.GetIDAndValidJWT>SpecialErr", nil)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
		}
		// Пакуем UserID в context
		ctx := context.WithValue(r.Context(), CtxUserID, UserID)
		next(w, r.WithContext(ctx))
	}
}
