package handler

import (
	"net/http"

	_ "github.com/lib/pq"

	"github.com/boginskiy/Clicki/internal/service"
)

type HandlerURL struct {
	Service service.CRUDer // CRUDer is the interface of business logic
}

func (h *HandlerURL) Get(res http.ResponseWriter, req *http.Request) {
	// Запуск бизнес логики сервиса 'Service'
	body, err := h.Service.Read(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", string(body))
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *HandlerURL) Post(res http.ResponseWriter, req *http.Request) {
	// Start of 'Service'
	body, err := h.Service.Create(req)

	// Check err
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	tmpHeader := h.Service.GetHeader()
	res.Header().Set("Content-Type", tmpHeader)
	res.WriteHeader(http.StatusCreated)
	res.Write(body)
}

func (h *HandlerURL) Check(res http.ResponseWriter, req *http.Request) {
	body, err := h.Service.CheckPing(req)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)
	res.Write(body)
}

func (h *HandlerURL) Set(res http.ResponseWriter, req *http.Request) {
	body, err := h.Service.SetBatch(req)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(body)
}
