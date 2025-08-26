package middleware

import (
	"net/http"
	"time"

	l "github.com/boginskiy/Clicki/internal/logger"
)

type Middleware struct {
	Logger l.Logger
}

func NewMiddleware(logger l.Logger) *Middleware {
	return &Middleware{Logger: logger}
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
