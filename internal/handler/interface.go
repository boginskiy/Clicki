package handler

import (
	mv "github.com/boginskiy/Clicki/internal/middleware"
	"github.com/go-chi/chi"
)

// Handler - .
// type Handler interface {

// 	Get(res http.ResponseWriter, req *http.Request)
// 	Post(res http.ResponseWriter, req *http.Request)
// 	Check(res http.ResponseWriter, req *http.Request)
// }

type Handler interface {
	RegisterRoutes(route chi.Router, mdlwere mv.Middleware)
}
