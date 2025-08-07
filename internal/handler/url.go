package handler

import (
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
		http.Error(res, ErrBodyReq, http.StatusBadRequest)
		return
	}

	// Генерируем ключ. Записываем originURL.
	imitationPath := h.encryptionLongURL()
	db.Store[imitationPath] = originURL

	// Host
	// host := req.Host
	// if config.ArgsCLI.IsCh {
	// 	host = h.ChangePort(req.Host, config.ArgsCLI.ResultPort)
	// }

	// Подготавливаем тело res. Формат 'http//localhost:8080/Jgd63Kd8'
	// imitationURL := fmt.Sprintf(
	// 	"%s://%s%s",
	// 	h.GetProtocol(req), host, req.URL.Path+imitationPath)

	imitationURL := config.ArgsCLI.ResultPort + imitationPath

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(imitationURL))
}
