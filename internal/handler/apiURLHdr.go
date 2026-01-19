package handler

import (
	"net/http"

	mv "github.com/boginskiy/Clicki/internal/middleware"
	p "github.com/boginskiy/Clicki/internal/protocol"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/go-chi/chi"
)

type APIURLHandlers struct {
	APIURLServ service.APIServicer
	APIDelServ service.DelServicer
	Protocol   p.Protocol
}

func NewAPIURLHandlers(
	apiURLServ service.APIServicer,
	apiDelServ service.DelServicer,
	prot p.Protocol) *APIURLHandlers {

	return &APIURLHandlers{
		APIURLServ: apiURLServ,
		APIDelServ: apiDelServ,
		Protocol:   prot,
	}
}

func (a *APIURLHandlers) RegisterRoutes(r chi.Router, mdlwere mv.Middleware) {
	r.Post("/shorten", mdlwere.Conveyor(http.HandlerFunc(a.Create)))
	r.Post("/shorten/batch", mdlwere.Conveyor(http.HandlerFunc(a.CreateSet)))
	r.Delete("/user/urls", mdlwere.Conveyor(http.HandlerFunc(a.DeleteSet)))
	r.Get("/user/urls", mdlwere.Conveyor(http.HandlerFunc(a.ReadSet)))
	r.Get("/internal/stats", mdlwere.WithTrustedSubnet(http.HandlerFunc(a.ShowStats)))
}

func (a *APIURLHandlers) ShowStats(w http.ResponseWriter, r *http.Request) {
	dataByte, err := a.APIURLServ.GetStats(r)
	if err != nil {
		http.Error(w, "message: not created", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dataByte)
}

func (a *APIURLHandlers) Create(w http.ResponseWriter, r *http.Request) {
	// Put in "Create" obj "Protocol" for processing "request" in "APIURLServ".
	dataByte, err := a.APIURLServ.Create(r.Context(), a.Protocol, r)

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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(dataByte)
}

func (a *APIURLHandlers) CreateSet(w http.ResponseWriter, r *http.Request) {
	body, err := a.APIURLServ.CreateSet(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func (a *APIURLHandlers) ReadSet(w http.ResponseWriter, r *http.Request) {
	data, err := a.APIURLServ.ReadSet(r.Context(), a.Protocol)
	dataByte, ok := data.([]byte)

	if err != nil || !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// if user does't have record.
	if len(dataByte) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dataByte)
}

func (a *APIURLHandlers) DeleteSet(w http.ResponseWriter, r *http.Request) {
	if _, err := a.APIDelServ.DeleteSet(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
