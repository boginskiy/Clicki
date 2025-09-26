package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/boginskiy/Clicki/internal/gzip"
	"github.com/boginskiy/Clicki/internal/logg"
	regstr "github.com/boginskiy/Clicki/internal/register"
	"github.com/boginskiy/Clicki/internal/repository"
)

type MvFunc func(http.HandlerFunc) http.HandlerFunc

type Middleware struct {
	Register regstr.Register
	RepoUser repository.RepositoryUser
	Logger   logg.Logger
}

func NewMiddleware(logger logg.Logger, register regstr.Register, repoUser repository.RepositoryUser) *Middleware {
	return &Middleware{
		Register: register,
		RepoUser: repoUser,
		Logger:   logger}
}

func (m *Middleware) Conveyor(next http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range []MvFunc{m.WithInfoLogger, m.WithGzip} { // m.WithAuth,
		next = middleware(next)
	}
	return next
}

// WithInfoLogger - логирование данных
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

// WithGzip - сжатие данных
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

// WithAuth - регистрация и аутентификация пользователя
func (m *Middleware) WithAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Идентификатор пользователя
		var UserID int

		body, cookie, err := m.Register.Registration(r, &UserID)
		if body == nil {
			body, cookie, err = m.Register.Authentication(r, &UserID)
		}

		// Пакуем идентификатор в context
		ctx := context.WithValue(r.Context(), CtxUserID, UserID)

		// Если какие-либо ошибки
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(body)
			return
		}

		// Если были сгенерированы куки
		if cookie != nil {
			http.SetCookie(w, cookie)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(body)
			return
		}

		// Если пользователь не зарегистрирован
		if body != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(body)
			return
		}

		// Если пользователь прошел регистрацию/аутентификацию, передаем управление конвейеру
		next(w, r.WithContext(ctx))
	}
}
