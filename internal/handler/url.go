package handler

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/service"
)

type HandlerURL struct {
	Service service.CRUDer   // CRUDer is the interface of business logic
	Kwargs  config.VarGetter // Kwargs is the args of command line interface
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
	body, err := h.Service.Create(req, h.Kwargs)

	// Check err
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Примитивная проверка, что перед нами Json, в зависимости от этого меняем тип 'Content-Type'
	tmpBody := []byte(h.Kwargs.GetBaseURL() + "/" + string(body))
	tmpHeader := "text/plain"

	if len(body) > 0 {
		switch body[0] {
		case '{', '[', '"':
			tmpHeader = "application/json"
			tmpBody = body
		}
	}
	res.Header().Set("Content-Type", tmpHeader)
	res.WriteHeader(http.StatusCreated)
	res.Write(tmpBody)
}

func (h *HandlerURL) Check(res http.ResponseWriter, req *http.Request) {
	body, err := h.Service.CheckPing(req)
	log.Println("body", body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)
	res.Write(body)
}
