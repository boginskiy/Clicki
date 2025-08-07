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

	// Генерируем ключ. Записываем originURL.
	imitationPath := h.encryptionLongURL()
	db.Store[imitationPath] = originURL

	// Параметры для сборки ответа
	typeProtocol := h.GetProtocol(req)

	host := req.Host
	if config.ArgsCLI.IsCh {
		host = h.ChangePort(req.Host, config.ArgsCLI.ResultPort)
	}

	path := req.URL.Path + imitationPath

	// Подготавливаем тело res. Формат 'http//localhost:8080/Jgd63Kd8'
	imitationURL := fmt.Sprintf(
		"%s://%s%s",
		typeProtocol, host, path)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	fmt.Fprintf(res, "%s", imitationURL)
}
