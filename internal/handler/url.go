package handler

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
)

const CAP = 10

type HandlerForURL struct {
	Store        map[string]string
	SetOfSymbols []rune
}

func NewHandlerForURL() *HandlerForURL {
	return &HandlerForURL{
		Store:        make(map[string]string, CAP),
		SetOfSymbols: Symbols,
	}
}

func (h *HandlerForURL) createImitationPath() string {
	slByte := make([]byte, 8)
	for {
		for i := range slByte {
			slByte[i] = byte(h.SetOfSymbols[rand.Intn(len(h.SetOfSymbols))])
		}
		// Проверка на уникальность
		if _, ok := h.Store[string(slByte)]; !ok {
			break
		}
	}
	return string(slByte)
}

// Костыль до использования иных библиотек с обработкой динамических url
func (h *HandlerForURL) checkUpPath(path string) bool {
	re := regexp.MustCompile(CheckPath)
	return re.MatchString(path)
}

func (h *HandlerForURL) checkUpBody(body string) bool {
	re := regexp.MustCompile(CheckDomain)
	return re.MatchString(body)
}

func (h *HandlerForURL) checkUpPathAndMethod(path, method string) bool {
	if h.checkUpPath(path) && method == "GET" {
		return true
	} else if path == "/" && method == "POST" {
		return true
	}
	return false
}

func (h *HandlerForURL) HandlerPostURL(req *http.Request) (string, error) {
	// Тащим строку из тела запроса
	tmpBody, err := io.ReadAll(req.Body)
	if err != nil {
		return "", errors.Join(errors.New(ErrBodyReq), err)
	}

	// Валидируем строку. Проверка, что строка является доменом сайта
	if !h.checkUpBody(string(tmpBody)) {
		return "", errors.New(ErrBodyReq)
	}

	// Записываем домен под сгенерированным ключом
	imitationPath := h.createImitationPath()
	h.Store[imitationPath] = string(tmpBody)
	return imitationPath, nil
}

func (h *HandlerForURL) HandlerGetURL(req *http.Request) (string, error) {
	// Достаем оригинальный URL
	tmpPath := req.URL.Path
	originPath, ok := h.Store[strings.TrimLeft(tmpPath, "/")]

	if !ok {
		return "", errors.New(ErrNotData)
	}
	return originPath, nil
}

func (h HandlerForURL) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Проверка валидности запроса и path
	if !h.checkUpPathAndMethod(req.URL.Path, req.Method) {
		http.Error(res, ErrPathAndMethod, http.StatusMethodNotAllowed)
		return
	}

	// Прокидываем на обработчики
	switch req.Method {

	case "POST":
		imitationPath, err := h.HandlerPostURL(req)

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
		originPath, err := h.HandlerGetURL(req)

		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		// Формируем тело ответа ...
		res.Header().Set("Location", originPath)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}
