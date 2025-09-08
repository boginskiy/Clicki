package router

import (
	"net/http"
	"net/http/pprof"

	c "github.com/boginskiy/Clicki/cmd/config"

	h "github.com/boginskiy/Clicki/internal/handler"
	m "github.com/boginskiy/Clicki/internal/middleware"
	s "github.com/boginskiy/Clicki/internal/service"
	"github.com/go-chi/chi"
)

func Router(kwargs c.VarGetter, mv m.Middlewarer, apiURL, shortuRL s.CRUDer) *chi.Mux {
	hAPIURL := h.HandlerURL{Service: apiURL, Kwargs: kwargs}
	hURL := h.HandlerURL{Service: shortuRL, Kwargs: kwargs}

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {

		// shortURL
		r.Post("/", mv.Conveyor(http.HandlerFunc(hURL.Post)))
		r.Get("/{id}", mv.Conveyor(http.HandlerFunc(hURL.Get)))
		r.Get("/ping", mv.WithInfoLogger(http.HandlerFunc(hURL.Check)))

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
