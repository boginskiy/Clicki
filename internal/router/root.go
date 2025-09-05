package router

import (
	"net/http"
	"net/http/pprof"

	c "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/db"
	h "github.com/boginskiy/Clicki/internal/handler"
	l "github.com/boginskiy/Clicki/internal/logger"
	m "github.com/boginskiy/Clicki/internal/middleware"
	p "github.com/boginskiy/Clicki/internal/preparation"
	s "github.com/boginskiy/Clicki/internal/service"
	v "github.com/boginskiy/Clicki/internal/validation"
	"github.com/go-chi/chi"
)

func Router(kwargs c.VarGetter, mv m.Middlewarer, db db.Storage, logger l.Logger) *chi.Mux {
	extraFuncer := p.NewExtraFunc() // extraFuncer - дополнительные возможности
	checker := v.NewChecker()       // checker - валидация данных

	APIShortURL := s.NewAPIShortURL(db, logger, checker, extraFuncer) // Service 'APIShortURL'
	shortURL := s.NewShortURL(db, logger, checker, extraFuncer)       // Service 'ShortURL'

	hURL := h.HandlerURL{Service: shortURL, Kwargs: kwargs}
	hAPIURL := h.HandlerURL{Service: APIShortURL, Kwargs: kwargs}

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {

		// shortURL
		r.Post("/", mv.Conveyor(http.HandlerFunc(hURL.Post)))
		r.Get("/{id}", mv.Conveyor(http.HandlerFunc(hURL.Get)))

		// APIShortURL
		r.Route("/api/", func(r chi.Router) {
			r.Post("/shorten", mv.Conveyor(http.HandlerFunc(hAPIURL.Post)))
		})

		// PProf
		r.Route("/debug/pprof/", func(r chi.Router) {
			r.Get("/", pprof.Index)
			r.Get("/cmdline", pprof.Cmdline)
			r.Get("/profile", pprof.Profile)
			r.Get("/symbol", pprof.Symbol)
			r.Get("/trace", pprof.Trace)
		})
	})
	return r
}
