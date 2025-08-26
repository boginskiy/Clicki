package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	c "github.com/boginskiy/Clicki/cmd/config"
	db "github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/logger"
	m "github.com/boginskiy/Clicki/internal/middleware"
	p "github.com/boginskiy/Clicki/internal/preparation"
	r "github.com/boginskiy/Clicki/internal/router"
	s "github.com/boginskiy/Clicki/internal/service"
	v "github.com/boginskiy/Clicki/internal/validation"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func RunRouter() *chi.Mux {
	kwargs := c.NewVariables()      // agrs - атрибуты командной строки
	extraFuncer := p.NewExtraFunc() // extraFuncer - дополнительные возможности
	checker := v.NewChecker()       // checker - валидация данных
	db := db.NewDBStore()           // db - слой базы данных 'DBStore'

	infoLog := logger.NewLogg("Test.log", "INFO")
	midWare := m.NewMiddleware(infoLog)

	// Заполняем базу данных тестовыми данными
	db.Store["DcKa7J8d"] = "https://translate.yandex.ru/"

	// shortingURL - слой с бизнес логикой сервиса 'ShorteningURL'
	shortingURL := s.NewShorteningURL(db, checker, extraFuncer, infoLog)
	//
	return r.Router(kwargs, midWare, shortingURL)
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

func TestRouter(t *testing.T) {
	ts := httptest.NewServer(RunRouter())
	defer ts.Close()

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
			res, _ := ExecuteRequest(t, ts, tt.methodReq, tt.urlReq, tt.bodyReq)
			defer res.Body.Close()
			assert.Equal(t, tt.statusRes, res.StatusCode)
			assert.Equal(t, tt.contentTypeRes, res.Header.Get(tt.contentViewRes))
		})
	}
}
