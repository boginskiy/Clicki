package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/boginskiy/Clicki/pkg"
)

const (
	CAP  = 10
	LONG = 8
)

type RootHandler struct {
	Store map[string]string
	pkg.Tools
}

func NewRootHandler() *RootHandler {
	return &RootHandler{
		Store: make(map[string]string, CAP),
	}
}

// TODO Это бизнес логика. Перенеси ее туда сам знаешь куда.
func (h *RootHandler) encryptionLongURL() (shortURL string) {
	for {
		// Вызов шифратора
		shortURL = pkg.Scramble(LONG)
		// Проверка на уникальность
		if _, ok := h.Store[shortURL]; !ok {
			break
		}
	}
	return shortURL
}

func (h *RootHandler) checkUpPathAndMethod(path, method string) bool {
	if h.CheckUpPath(path) && method == "GET" {
		return true
	} else if path == "/" && method == "POST" {
		return true
	}
	return false
}

func (h *RootHandler) HandlerPost(req *http.Request) (string, error) {
	// Тащим строку из тела запроса
	tmpBody, err := io.ReadAll(req.Body)
	if err != nil {
		return "", errors.Join(errors.New(ErrBodyReq), err)
	}

	// Валидируем строку. Проверка, что строка является доменом сайта
	if !h.CheckUpBody(string(tmpBody)) {
		return "", errors.New(ErrBodyReq)
	}

	// Записываем домен под сгенерированным ключом
	imitationPath := h.encryptionLongURL()
	h.Store[imitationPath] = string(tmpBody)
	return imitationPath, nil
}

func (h *RootHandler) HandlerGet(req *http.Request) (string, error) {
	// Достаем оригинальный URL
	tmpPath := req.URL.Path
	originPath, ok := h.Store[strings.TrimLeft(tmpPath, "/")]

	if !ok {
		return "", errors.New(ErrNotData)
	}
	return originPath, nil
}

func (h RootHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Проверка валидности запроса и path
	if !h.checkUpPathAndMethod(req.URL.Path, req.Method) {
		http.Error(res, ErrPathAndMethod, http.StatusMethodNotAllowed)
		return
	}

	// Прокидываем на обработчики
	switch req.Method {

	case "POST":
		imitationPath, err := h.HandlerPost(req)

		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		// Формируем тело ответа ...
		resBody := fmt.Sprintf("http://%s/%s", req.Host, imitationPath)

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(resBody))

	case "GET":
		originPath, err := h.HandlerGet(req)

		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		// Формируем тело ответа ...
		res.Header().Set("Location", originPath)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}
