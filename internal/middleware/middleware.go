package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/boginskiy/Clicki/internal/gzip"
	l "github.com/boginskiy/Clicki/internal/logger"
)

type MvFunc func(http.HandlerFunc) http.HandlerFunc

type Middleware struct {
	Logger l.Logger
}

func NewMiddleware(logger l.Logger) *Middleware {
	return &Middleware{Logger: logger}
}

func (m *Middleware) Conveyor(next http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range []MvFunc{m.WithInfoLogger, m.WithGzip} {
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
		extW := NewExtRespWrtr(w)
		next(extW, r)

		duration := time.Since(start)

		m.Logger.RaiseInfo(
			l.DataReqResInfo,
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
