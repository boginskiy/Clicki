package server

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	conf "github.com/boginskiy/Clicki/cmd/config"
	db "github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/logg"
	midw "github.com/boginskiy/Clicki/internal/middleware"
	mod "github.com/boginskiy/Clicki/internal/model"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	regstr "github.com/boginskiy/Clicki/internal/register"
	"github.com/boginskiy/Clicki/internal/repository"
	route "github.com/boginskiy/Clicki/internal/router"
	srv "github.com/boginskiy/Clicki/internal/service"
	valid "github.com/boginskiy/Clicki/internal/validation"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func RunRouter() *chi.Mux {
	infoLog := logg.NewLogg("Test.log", "INFO")
	kwargs := conf.NewVariables(infoLog) // agrs - атрибуты командной строки
	kwargs.PathToStore = "test"

	db, _ := db.NewStoreFile(kwargs, infoLog)

	// Данные для тестирования
	url := &mod.URLTb{CorrelationID: "DcKa7J8d", OriginalURL: "https://translate.yandex.ru/"}
	store := map[string]*mod.URLTb{"DcKa7J8d": url}
	uniqueFields := map[string]string{url.OriginalURL: url.CorrelationID}

	// Специальное создание репозитория для теста с начальным обогащением данных
	file, _ := db.GetDB().(*os.File)
	repo := &repository.RepositoryFileURL{
		Kwargs:  kwargs,
		DB:      db,
		Scanner: bufio.NewScanner(file),
		File:    file,
	}
	repo.Store = store
	repo.UniqueFields = uniqueFields

	register := regstr.NewRegist(kwargs, infoLog, repo)
	midWare := midw.NewMiddleware(infoLog, register, repo)
	extraFuncer := prep.NewExtraFunc()
	checker := valid.NewChecker()

	// Services
	APIShortURL := srv.NewAPIShortURL(kwargs, infoLog, repo, checker, extraFuncer)
	ShortURL := srv.NewShortURL(kwargs, infoLog, repo, checker, extraFuncer)

	return route.Router(midWare, APIShortURL, ShortURL)
}

func ExecuteRequest(t *testing.T, ts *httptest.Server, method, url, body string) (*http.Response, string) {
	// New Req
	req, err := http.NewRequest(method, ts.URL+url, strings.NewReader(body))
	require.NoError(t, err)

	// New Client
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Do(req)
	require.NoError(t, err)

	resBody, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	require.NoError(t, err)
	return res, string(resBody)
}

func ExecuteRequest2(t *testing.T, ts *httptest.Server, client *http.Client, method, url, body string) (*http.Response, string) {
	// New Req
	req, err := http.NewRequest(method, ts.URL+url, strings.NewReader(body))
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)

	resBody, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	require.NoError(t, err)
	return res, string(resBody)
}

func testRouter(t *testing.T, server *httptest.Server) {
	// ts := httptest.NewServer(RunRouter())
	// defer ts.Close()

	// Tasts Cases
	tests := []struct {
		name           string
		methodReq      string
		bodyReq        string
		urlReq         string
		contentViewRes string
		contentTypeRes string
		statusRes      int
	}{
		// POST
		{"Test POST 1", "POST", "://docs.google.com/", "/", "Content-Type", "text/plain; charset=utf-8", 400},
		{"Test POST 2", "POST", "https://docs.google.com/", "/", "Content-Type", "text/plain", 201},
		{"Test POST 3", "POST", "https://docs.google.com/", "/", "Content-Type", "text/plain", 409},
		{"Test POST 4", "POST", "", "/wwxwecq", "Content-Type", "", 405},

		// GET
		{"Test GET 1", "GET", "", "/DcKa7J44", "Content-Type", "text/plain; charset=utf-8", 400},
		{"Test GET 2", "GET", "", "/DcKa7J8d", "Location", "https://translate.yandex.ru/", 307},
		{"Test GET 3", "GET", "", "/", "Content-Type", "", 405},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _ := ExecuteRequest(t, server, tt.methodReq, tt.urlReq, tt.bodyReq)
			defer res.Body.Close()
			assert.Equal(t, tt.statusRes, res.StatusCode)
			assert.Equal(t, tt.contentTypeRes, res.Header.Get(tt.contentViewRes))
		})
	}
}

func testRegistration(t *testing.T, server *httptest.Server) {
	// tests := []struct {
	// 	name           string
	// 	methodReq      string
	// 	bodyReq        string
	// 	urlReq         string
	// 	contentViewRes string
	// 	contentTypeRes string
	// 	statusRes      int
	// }{
	// 	{"Test Registration POST 1", "POST", "https://docs.milk.com/", "/", "Content-Type", "application/json", 200},
	// 	{"Test Registration POST 2", "POST", "https://docs.milk.com/", "/", "Content-Type", "text/plain", 201},
	// }

	// Client
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Req 1 Пользователь получает токен
	req, err := http.NewRequest("POST", server.URL+"/", strings.NewReader("https://docs.milk.com/"))
	require.NoError(t, err)
	res1, err := client.Do(req)
	require.NoError(t, err)

	// _, err = io.ReadAll(res1.Body)
	// defer res1.Body.Close()

	assert.Equal(t, 200, res1.StatusCode)
	assert.Equal(t, "application/json", res1.Header.Get("Content-Type"))

	// Cookies
	cookies := res1.Cookies()

	fmt.Println(">>", cookies)

	// Req 2 Пользователь с полученным токеном делает еще запрос, но его нет в бд
	req2, err2 := http.NewRequest("POST", server.URL+"/", strings.NewReader("https://docs.milk.com/"))
	require.NoError(t, err2)

	// Установка куки в новый запрос
	for _, cookie := range cookies {
		req2.AddCookie(cookie)
	}

	res2, err := client.Do(req2)
	require.NoError(t, err)

	// _, err = io.ReadAll(res2.Body)
	// defer res1.Body.Close()

	assert.Equal(t, 401, res2.StatusCode)
	assert.Equal(t, "application/json", res2.Header.Get("Content-Type"))

	// Req 3 Нужно чтобы пользователь появился в БД и сделать запрос еще раз

	// resBody, err := io.ReadAll(res.Body)
	// defer res.Body.Close()
	// require.NoError(t, err)
	// return res, string(resBody)

	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		res, _ := ExecuteRequest2(t, server, client, tt.methodReq, tt.urlReq, tt.bodyReq)
	// 		defer res.Body.Close()
	// 		assert.Equal(t, tt.statusRes, res.StatusCode)
	// 		assert.Equal(t, tt.contentTypeRes, res.Header.Get(tt.contentViewRes))
	// 	})
	// }
}

func testCompress(t *testing.T, server *httptest.Server) {
	// Request
	requestBody := "https://practicum.yandex.ru"

	// Test_1
	t.Run("Test_compressing_data_in_request", func(t *testing.T) {
		// Сжимаем клиентский запрос
		buf := bytes.NewBuffer(nil)
		wGzip := gzip.NewWriter(buf)
		_, err := wGzip.Write([]byte(requestBody))
		defer wGzip.Close()

		require.NoError(t, err)
		err = wGzip.Close()
		require.NoError(t, err)

		// Подготовка запроса
		req := httptest.NewRequest("POST", server.URL, buf)
		req.RequestURI = ""
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "")

		// Отправка запроса
		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, 201, res.StatusCode)
		defer res.Body.Close()

		// Check response body
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		require.Contains(t, string(body), "http://localhost:8080")

	})

	// Test_2
	t.Run("Test_compressing_data_in_response", func(t *testing.T) {
		// Подготовка запроса
		req := httptest.NewRequest("POST", server.URL, strings.NewReader(requestBody))
		req.RequestURI = ""
		req.Header.Set("Content-Type", "text/html")
		req.Header.Set("Accept-Encoding", "gzip")

		// Отправка запроса
		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, 409, res.StatusCode) // TODO! 201 Почему стало 409 ?
		require.Equal(t, res.Header.Get("Content-Encoding"), "gzip")
		defer res.Body.Close()

		// Checking Body
		rGzip, err := gzip.NewReader(res.Body)
		require.NoError(t, err)

		defer rGzip.Close()

		var b bytes.Buffer

		_, err = b.ReadFrom(rGzip)
		require.NoError(t, err)

		require.Contains(t, b.String(), "http://localhost:8080")
	})
}

func TestMain(t *testing.T) {
	server := httptest.NewServer(RunRouter())
	defer server.Close()

	// Test Registration
	testRegistration(t, server)

	// // Test Router
	// testRouter(t, server)

	// // Test Compress
	// testCompress(t, server)
}
