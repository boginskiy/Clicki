package router

import (
	"net/http"
	"net/http/pprof"

	"github.com/boginskiy/Clicki/internal/handler"
	midw "github.com/boginskiy/Clicki/internal/middleware"
	srv "github.com/boginskiy/Clicki/internal/service"
	"github.com/go-chi/chi"
)

func Router(mv midw.Middlewarer, apiURL, shortuRL srv.CrudSrver, apiDelMess srv.DelSrver) *chi.Mux {
	hAPIURL := handler.HandlerURL{CrudSrver: apiURL, DelSrver: apiDelMess}
	hURL := handler.HandlerURL{CrudSrver: shortuRL, DelSrver: nil}

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {

		// shortURL
		r.Post("/", mv.Conveyor(http.HandlerFunc(hURL.CreateURL)))
		r.Get("/{id}", mv.Conveyor(http.HandlerFunc(hURL.ReadURL)))
		r.Get("/ping", mv.WithInfoLogger(http.HandlerFunc(hURL.CheckDB)))

		// APIShortURL
		r.Route("/api/", func(r chi.Router) {
			r.Post("/shorten", mv.Conveyor(http.HandlerFunc(hAPIURL.CreateURL)))
			r.Post("/shorten/batch", mv.Conveyor(http.HandlerFunc(hAPIURL.CreateSetURL)))
			r.Get("/user/urls", mv.Conveyor(http.HandlerFunc(hAPIURL.ReadSetUserURL)))
			r.Delete("/user/urls", mv.Conveyor(http.HandlerFunc(hAPIURL.DeleteSetUserURL)))
		})

		// PProf
		r.Route("/debug/pprof/", func(r chi.Router) {
			r.Get("/", pprof.Index)
			r.Get("/cmdline", pprof.Cmdline)
			r.Get("/profile", pprof.Profile)
			r.Get("/symbol", pprof.Symbol)
			r.Get("/trace", pprof.Trace)
			r.Get("/heap", func(w http.ResponseWriter, r *http.Request) {
				pprof.Handler("heap").ServeHTTP(w, r)
			})
		})
	})
	return r
}
