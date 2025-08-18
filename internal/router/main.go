package router

import (
	"net/http"

	c "github.com/boginskiy/Clicki/cmd/config"
	h "github.com/boginskiy/Clicki/internal/handler"
	s "github.com/boginskiy/Clicki/internal/service"
	"github.com/go-chi/chi"
)

func Router(shortingURL s.ShortenerURL, argsCLI *c.ArgumentsCLI) *chi.Mux {
	h := h.RootHandler{ShortingURL: shortingURL, ArgsCLI: argsCLI}
	r := chi.NewRouter()

	// Деревце routes
	r.Route("/", func(r chi.Router) {
		r.Post("/", http.HandlerFunc(h.PostURL))
		r.Get("/{id}", http.HandlerFunc(h.GetURL))
	})
	return r
}
