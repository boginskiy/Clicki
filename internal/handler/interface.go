package handler

import (
	mv "github.com/boginskiy/Clicki/internal/middleware"
	"github.com/go-chi/chi"
)

type Handler interface {
	RegisterRoutes(route chi.Router, mdlwere mv.Middleware)
}
