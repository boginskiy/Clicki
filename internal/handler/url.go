package handler

import (
	"net/http"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/service"
)

type RootHandler struct {
	ShortingURL service.ShortenerURL // ShortingURL is the interface of business logic
	Kwargs      config.Variabler     // ArgsCLI is the args of command line interface
}

func (h *RootHandler) GetURL(res http.ResponseWriter, req *http.Request) {
	// Запуск бизнес логики сервиса 'ShorteningURL'
	err := h.ShortingURL.Execute(req)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	// Сборка ответа
	res.Header().Set("Location", h.ShortingURL.GetOriginURL())
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *RootHandler) PostURL(res http.ResponseWriter, req *http.Request) {
	// Запуск бизнес логики сервиса 'ShorteningURL'
	err := h.ShortingURL.Execute(req)

	// Проверка ошибок
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Сборка ответа
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	// TODO. h.ArgsCLI.ResultPort
	res.Write([]byte(h.Kwargs.GetBaseURL() + "/" + h.ShortingURL.GetImitationPath()))
}
