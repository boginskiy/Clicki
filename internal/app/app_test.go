package app_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/audit"
	"github.com/boginskiy/Clicki/internal/auth"
	"github.com/boginskiy/Clicki/internal/database"
	"github.com/boginskiy/Clicki/internal/handler"
	"github.com/boginskiy/Clicki/internal/logg"
	mv "github.com/boginskiy/Clicki/internal/middleware"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/protocol"
	"github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/internal/router"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/internal/tester/tfunc"
	"github.com/boginskiy/Clicki/internal/tester/tserv"
	"github.com/boginskiy/Clicki/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	// Logger & Config
	pathToStore := "test_store"
	pathToLogg := "test.log"

	logg := logg.NewLogg(pathToLogg, "INFO")
	config := tserv.InitConfig()
	config.PathToStore = pathToStore

	server := httptest.NewServer(RunRouter(config, logg))

	testRouter(t, server)   // Test Router
	testCompress(t, server) // Test Compress

	defer tfunc.DeleteTestFiles(pathToStore, pathToLogg)
	defer server.Close()
}

func RunRouter(config config.Config, logg logg.Logger) http.Handler {
	// Database & Repo.
	database, _ := database.NewStoreFile(config, logg)
	repo := repository.NewMainRepoFile(config, logg, database)

	// Auth & middleware.
	auther := auth.NewAuth(config, logg, repo)
	midWare := mv.NewMdlwere(config, logg, auther)

	// Some function.
	fancer := prep.NewFunctions(logg)
	checker := validation.NewChecker()

	// Ctx.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обогощение репозитория данными.
	tfunc.WriteRecord(repo)

	// Publisher
	var sub1 = audit.NewFileReceiver(logg, config.GetAuditFile(), 1)
	var sub2 = audit.NewServerReceiver(logg, config.GetAuditURL(), 2)
	var publisher = audit.NewPublish(sub1, sub2)

	// Services
	APIURLServ := service.NewAPIURLServ(config, logg, repo, checker, fancer, publisher)
	URLServ := service.NewURLServ(config, logg, repo, checker, fancer, publisher)
	APIDelServ := service.NewAPIDelServ(ctx, config, logg, repo)

	// Protocol
	protHTTP := protocol.NewProtocolHTTP()

	// Handler
	APIURLHdler := handler.NewAPIURLHandlers(APIURLServ, APIDelServ, protHTTP)
	URLHdler := handler.NewURLHandlers(URLServ)
	PprofHdler := handler.NewPprofHandlers()

	// Router
	return router.NewRoute(URLHdler, APIURLHdler, PprofHdler).Run(midWare)
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
		{"Test GET 2", "GET", "", "/wrs4db6j", "Location", "https://practicum.yandex.ru/", 307},
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
	t.Run("test compressing data in request", func(t *testing.T) {
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
	t.Run("test compressing data in response", func(t *testing.T) {
		// Подготовка запроса
		req := httptest.NewRequest("POST", server.URL, strings.NewReader(requestBody))
		req.RequestURI = ""
		req.Header.Set("Content-Type", "text/html")
		req.Header.Set("Accept-Encoding", "gzip")

		// Отправка запроса
		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, 409, res.StatusCode)
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
