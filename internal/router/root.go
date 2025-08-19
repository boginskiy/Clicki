package router

import (
	"net/http"
	"net/http/pprof"

	c "github.com/boginskiy/Clicki/cmd/config"
	h "github.com/boginskiy/Clicki/internal/handler"
	s "github.com/boginskiy/Clicki/internal/service"
	"github.com/go-chi/chi"
)

// PProf
// $go tool pprof http://localhost:8080/debug/pprof/profile
// $ab -k -c 10 -n 100000 "http://127.0.0.1:8080/time"

func Router(shortingURL s.ShortenerURL, argsCLI *c.ArgumentsCLI) *chi.Mux {
	h := h.RootHandler{ShortingURL: shortingURL, ArgsCLI: argsCLI}
	r := chi.NewRouter()

	// Tree routes
	r.Route("/", func(r chi.Router) {
		r.Post("/", http.HandlerFunc(h.PostURL))
		r.Get("/{id}", http.HandlerFunc(h.GetURL))

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
