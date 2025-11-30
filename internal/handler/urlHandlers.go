package handler

import (
	"net/http"

	mv "github.com/boginskiy/Clicki/internal/middleware"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/go-chi/chi"
)

type URLHandlers struct {
	URLServ service.Servicer
}

func NewURLHandlers(urlServ service.Servicer) *URLHandlers {
	return &URLHandlers{
		URLServ: urlServ,
	}
}

func (u *URLHandlers) RegisterRoutes(r chi.Router, mdlwere mv.Middleware) {
	r.Post("/", mdlwere.Conveyor(http.HandlerFunc(u.Create)))
	r.Get("/{id}", mdlwere.Conveyor(http.HandlerFunc(u.Read)))
	r.Get("/ping", mdlwere.WithLogg(http.HandlerFunc(u.CheckDB)))
}

func (u *URLHandlers) Create(w http.ResponseWriter, r *http.Request) {
	dataByte, err := u.URLServ.Create(r)
	status := http.StatusCreated

	// Processing critical errors.
	if err != nil && len(dataByte) == 0 {
		http.Error(w, "message: not created", http.StatusBadRequest)
		return
	}
	// Processing uncritical errors.
	if err != nil && len(dataByte) > 0 {
		status = http.StatusConflict
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write(dataByte)
}

func (u *URLHandlers) Read(w http.ResponseWriter, r *http.Request) {
	dataByte, err := u.URLServ.Read(r)

	if err == service.ErrReadRecord {
		w.WriteHeader(http.StatusGone)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", string(dataByte))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (u *URLHandlers) CheckDB(w http.ResponseWriter, r *http.Request) {
	dataByte, err := u.URLServ.CheckDB(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(dataByte)
}
