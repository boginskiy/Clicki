package server

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	c "github.com/boginskiy/Clicki/cmd/config"
	db "github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/logger"
	m "github.com/boginskiy/Clicki/internal/middleware"
	"github.com/boginskiy/Clicki/internal/model"
	p "github.com/boginskiy/Clicki/internal/preparation"
	r "github.com/boginskiy/Clicki/internal/router"
	s "github.com/boginskiy/Clicki/internal/service"
	v "github.com/boginskiy/Clicki/internal/validation"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func RunRouter() *chi.Mux {
	infoLog := logger.NewLogg("Test.log", "INFO")
	kwargs := c.NewVariables(infoLog) // agrs - атрибуты командной строки
	kwargs.PathToStore = "test"

	repo, _ := db.NewStoreFile(kwargs, infoLog)
	midWare := m.NewMiddleware(infoLog)
	extraFuncer := p.NewExtraFunc()
	checker := v.NewChecker()

	// Services
	APIShortURL := s.NewAPIShortURL(repo, infoLog, checker, extraFuncer)
	ShortURL := s.NewShortURL(repo, infoLog, checker, extraFuncer)

	// type URLFile struct {
	// 	UUID        int    `json:"uuid"`
	// 	ShortURL    string `json:"short_url"`
	// 	OriginalURL string `json:"original_url"`
	// }

	// Заполняем базу данных тестовыми данными
	url := &model.URLFile{ShortURL: "DcKa7J8d", OriginalURL: "https://translate.yandex.ru/"}
	repo.Store["DcKa7J8d"] = url

	return r.Router(kwargs, midWare, APIShortURL, ShortURL)
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
		{"Test POST 3", "POST", "", "/wwxwecq", "Content-Type", "", 405},

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
		require.Equal(t, 201, res.StatusCode)
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

	// Test Router
	testRouter(t, server)

	// Test Compress
	testCompress(t, server)
}
