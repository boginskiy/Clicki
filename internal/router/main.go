package router

import (
	h "github.com/boginskiy/Clicki/internal/handler"
	"github.com/go-chi/chi"
)

func Router() *chi.Mux {
	h := h.RootHandler{}
	r := chi.NewRouter()

	// Деревце routes
	r.Route("/", func(r chi.Router) {
		r.Post("/", h.PostURL)
		r.Get("/{id}", h.GetURL)
	})
	return r
}
