package handler

import (
	"net/http"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/pkg/tools"
)

type RootHandler struct {
	ShortingURL service.ShortenerURL // Зависимость для ShortenerURL на interface
	ArgsCLI     *config.ArgumentsCLI // Аргументы командной строки
	tools.Tools                      // Вспомогательные инструменты
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
	res.Write([]byte(h.ArgsCLI.ResultPort + "/" + h.ShortingURL.GetImitationPath()))
}

// TODO
// На каналы перевести атрибуты структуры
// Тестирование
