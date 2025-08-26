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

func Router(kwargs c.VarGetter, mv m.Middlewarer, shortingURL s.ShortenerURL) *chi.Mux {
	h := h.RootHandler{ShortingURL: shortingURL, Kwargs: kwargs}
	r := chi.NewRouter()

	// Tree routes
	r.Route("/", func(r chi.Router) {
		r.Post("/", mv.WithInfoLogger(http.HandlerFunc(h.PostURL)))
		r.Get("/{id}", mv.WithInfoLogger(http.HandlerFunc(h.GetURL)))

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
