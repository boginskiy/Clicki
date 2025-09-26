package handler

import (
	"net/http"

	_ "github.com/lib/pq"

	srv "github.com/boginskiy/Clicki/internal/service"
)

type HandlerURL struct {
	Service srv.CRUDer
}

func (h *HandlerURL) ReadURL(res http.ResponseWriter, req *http.Request) {
	// Запуск бизнес логики сервиса 'Service'
	body, err := h.Service.ReadURL(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", string(body))
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *HandlerURL) CreateURL(res http.ResponseWriter, req *http.Request) {
	// Start of 'Service'
	body, err := h.Service.CreateURL(req)
	status := http.StatusCreated

	// Обработка критичных ошибок
	if err != nil && len(body) == 0 {
		http.Error(res, "message: not created", http.StatusBadRequest)
		return
	}

	// Обработка не критичных ошибок
	if err != nil && len(body) > 0 {
		status = http.StatusConflict
	}

	res.Header().Set("Content-Type", h.Service.GetHeader())
	res.WriteHeader(status)
	res.Write(body)
}

func (h *HandlerURL) CheckDB(res http.ResponseWriter, req *http.Request) {
	body, err := h.Service.CheckDB(req)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)
	res.Write(body)
}

func (h *HandlerURL) CreateSetURL(res http.ResponseWriter, req *http.Request) {
	body, err := h.Service.CreateSetURL(req)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write(body)
}

func (h *HandlerURL) ReadSetUserURL(res http.ResponseWriter, req *http.Request) {
	body, err := h.Service.ReadSetUserURL(req)

	// Ошибки
	// TODO! Надо бы привести все к одному виду/ Конечно ни факт что тесты пройдут.
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// У пользователя нет записей
	if len(body) == 0 {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusNoContent)
		res.Write(MessNoContent)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(body)
}
