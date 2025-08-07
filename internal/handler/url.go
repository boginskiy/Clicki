package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/pkg"
	"github.com/boginskiy/Clicki/pkg/tools"
)

const LONG = 8

type RootHandler struct {
	tools.Tools
}

// TODO - уберу
func (h *RootHandler) encryptionLongURL() (imitationURL string) {
	for {
		// Вызов шифратора
		imitationURL = pkg.Scramble(LONG)
		// Проверка на уникальность
		if _, ok := db.Store[imitationURL]; !ok {
			break
		}
	}
	return imitationURL
}

// TODO - уберу
func (h *RootHandler) checkUpPathAndMethod(path, method string) bool {
	if h.CheckUpPath(path) && method == "GET" {
		return true
	} else if path == "/" && method == "POST" {
		return true
	}
	return false
}

func (h *RootHandler) GetURL(res http.ResponseWriter, req *http.Request) {
	// Костыль для прохождения тестов. Уберу
	if !h.checkUpPathAndMethod(req.URL.Path, req.Method) {
		http.Error(res, ErrPathAndMethod, http.StatusMethodNotAllowed)
		return
	}

	// TODO. Оставил для прохождения локальных тестов
	tmpPath := req.URL.Path

	// // Достаем параметр id
	// imitationURL := chi.URLParam(req, "id")

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
	// Костыль для прохождения тестов. Уберу
	if !h.checkUpPathAndMethod(req.URL.Path, req.Method) {
		http.Error(res, ErrPathAndMethod, http.StatusMethodNotAllowed)
		return
	}

	// Вынимаем тело запроса
	originURL, err := h.TakeAllBodyFromReq(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Валидируем URL. Проверка регуляркой, что строка является доменом сайта
	if !h.CheckUpURL(originURL) || originURL == "" {
		http.Error(res, ErrBodyReq, http.StatusBadRequest)
		return
	}

	// Генерируем ключ. Записываем домен.
	imitationURL := h.encryptionLongURL()
	db.Store[imitationURL] = originURL

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)

	// Изменение порта при необходимости
	if config.ArgsCLI.IsCh {
		h.ChangePort(req, config.ArgsCLI.ResultPort)
	}

	// http//localhost:8080/Jgd63Kd8
	fmt.Fprintf(res, "%s://%s%s%s", h.GetProtocol(req), req.Host, req.URL.Path, imitationURL)
}
