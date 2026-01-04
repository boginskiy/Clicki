package handler

import (
	"net/http"
	"net/http/pprof"

	mv "github.com/boginskiy/Clicki/internal/middleware"
	"github.com/go-chi/chi"
)

type PprofHandlers struct {
}

func NewPprofHandlers() *PprofHandlers {
	return &PprofHandlers{}
}

func (p *PprofHandlers) RegisterRoutes(r chi.Router, _ mv.Middleware) {
	r.Get("/", pprof.Index)
	r.Get("/cmdline", pprof.Cmdline)
	r.Get("/profile", pprof.Profile)
	r.Get("/symbol", pprof.Symbol)
	r.Get("/trace", pprof.Trace)
	r.Get("/heap", func(
		w http.ResponseWriter, r *http.Request) {
		pprof.Handler("heap").ServeHTTP(w, r)
	})
}
