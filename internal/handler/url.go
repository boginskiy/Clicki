package handler

import (
	"net/http"
	"strings"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/pkg/tools"
)

type RootHandler struct {
	tools.Tools
}

func (h *RootHandler) GetURL(res http.ResponseWriter, req *http.Request) {
	// TODO. Оставил для прохождения локальных тестов
	tmpPath := req.URL.Path

	// // Достаем параметр id
	// imitationPath := chi.URLParam(req, "id")

	// Достаем origin URL
	originURL, ok := db.Store[strings.TrimLeft(tmpPath, "/")]
	if !ok {
		http.Error(res, "data is not available", http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", originURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *RootHandler) PostURL(res http.ResponseWriter, req *http.Request) {
	// Вынимаем тело запроса
	originURL, err := h.TakeAllBodyFromReq(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Валидируем URL. Проверка регуляркой, что строка является доменом сайта
	if !h.CheckUpURL(originURL) || originURL == "" {
		http.Error(res, "data not available or invalid", http.StatusBadRequest)
		return
	}

	// Генерируем ключ. Записываем originURL.
	imitationPath := service.ShortenerURL.EncryptionLongURL()
	db.Store[imitationPath] = originURL

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(config.ArgsCLI.ResultPort + "/" + imitationPath))
}
