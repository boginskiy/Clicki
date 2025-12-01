package router

import (
	"net/http"

	"github.com/boginskiy/Clicki/internal/handler"
	mv "github.com/boginskiy/Clicki/internal/middleware"
	"github.com/go-chi/chi"
)

// Router - interface
type Router interface {
	Run(mdlwere mv.Middleware) (handler http.Handler)
}

// Route - is struct about Router
type Route struct {
	R              *chi.Mux
	APIURLHandlers handler.Handler
	URLHandlers    handler.Handler
	PprofHandlers  handler.Handler
}

func NewRoute(urlHdler, apiURLHdler, pprofHdler handler.Handler) *Route {
	return &Route{
		R:              chi.NewRouter(),
		APIURLHandlers: apiURLHdler,
		PprofHandlers:  pprofHdler,
		URLHandlers:    urlHdler,
	}
}

func (r *Route) Run(mdlwere mv.Middleware) http.Handler {
	r.R.Route("/", func(route chi.Router) {
		r.URLHandlers.RegisterRoutes(route, mdlwere) // RLService.

		r.R.Route("/api/", func(route chi.Router) {
			r.APIURLHandlers.RegisterRoutes(route, mdlwere) // APIURLService.
		})

		r.R.Route("/debug/pprof/", func(route chi.Router) {
			r.PprofHandlers.RegisterRoutes(route, mdlwere) // PProf.
		})
	})
	return r.R
}
