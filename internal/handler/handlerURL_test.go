package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/audit"
	"github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/handler"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/model"
	"github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/internal/validation"
)

func TestHandlerURL(t *testing.T) {
	// Init
	logger := InitLogg("../tester/logs/TestHandlerURL.log", "INFO")
	kwargs := InitKwargs()
	db := InitDB()

	servShortURL := Init(logger, kwargs, db)

	// Testing
	testPostURL(t, servShortURL)
	testGetURL(t, servShortURL)
}

func Init(logger logg.Logger, kwargs config.VarGetter, db db.DBer) *service.ShortURL {
	var repo, _ = repository.NewRepositoryMapURL(kwargs, db)

	var sub1 = audit.NewFileReceiver(logger, kwargs.GetAuditFile(), 1)
	var sub2 = audit.NewServerReceiver(logger, kwargs.GetAuditURL(), 2)
	var publisher = audit.NewPublish(sub1, sub2)

	var core = service.NewCoreService(kwargs, logger, repo, publisher)
	var extraFuncer = preparation.NewExtraFunc()
	var checker = validation.NewChecker()

	return service.NewShortURL(core, repo, checker, extraFuncer)
}

func InitKwargs() config.VarGetter {
	return &config.Variables{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8081",
	}
}

func InitLogg(name, mod string) logg.Logger {
	return logg.NewLogg(name, mod)
}

func InitDB() db.DBer {
	return &db.StoreMap{Store: map[string]*model.URLTb{
		"H3HIkks3": {
			ID:            0,
			OriginalURL:   "https://practicum.yandex.ru/",
			ShortURL:      "short_url",
			CorrelationID: "H3HIkks3",
		}}}
}

// TestHandlerURL check only POST request
func testPostURL(t *testing.T, shURL *service.ShortURL) {
	type req struct {
		url  string
		body string
		host string
	}
	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name string
		want want
		req  req
	}{
		{
			name: "Test POST positive",
			want: want{
				contentType: "text/plain",
				statusCode:  201,
			},
			req: req{
				url:  "/",
				body: "https://practicum.yandex.ru/",
				host: "localhost:8080",
			},
		},

		{
			name: "Test POST negative",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
			},
			req: req{
				url:  "/",
				body: "jo55jt45oJOJOJPJOJJPWJP34O53R/",
				host: "localhost:8080",
			},
		},
	}

	//
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			// Request
			request := httptest.NewRequest(http.MethodPost, tt.req.url, strings.NewReader(tt.req.body))
			request.Host = tt.req.host
			// Recorder
			response := httptest.NewRecorder()
			// Handler
			h := handler.HandlerURL{CrudSrver: shURL, DelSrver: nil}
			h.CreateURL(response, request)

			// Check >>

			// StatusCode
			if response.Code != tt.want.statusCode {
				t.Errorf("%s:\n\texpected: %v\n\tactual: %v", tt.name, tt.want.statusCode, response.Code)
			}

			// Content-Type
			if response.Header().Get("Content-Type") != tt.want.contentType {
				t.Errorf("%s:\n\texpected: %v\n\tactual: %v", tt.name, tt.want.contentType, response.Header().Get("Content-Type"))
			}

			// Short URL
			if response.Code == 200 {
				tmpSl := strings.Split(response.Body.String(), "/")
				shortURL := tmpSl[len(tmpSl)-1]

				if len(shortURL) != service.LONG {
					t.Errorf("%s:\n\texpected: %v\n\tactual: %v", tt.name, len(shortURL), service.LONG)
				}
			}
			// <<
		})
	}
}

// TestHandlerURL2 check only GET request
func testGetURL(t *testing.T, shURL *service.ShortURL) {
	type req struct {
		url string
	}
	type want struct {
		contentType string
		location    string
		statusCode  int
	}

	tests := []struct {
		name  string
		want  want
		req   req
		store map[string]*model.URLTb
	}{
		{
			name: "Test GET positive",
			want: want{
				contentType: "",
				statusCode:  307,
				location:    "https://practicum.yandex.ru/",
			},
			req: req{
				url: "/H3HIkks3",
			},
			// store: map[string]*model.URLTb{
			// 	"H3HIkks3": {
			// 		ID:            0,
			// 		OriginalURL:   "https://practicum.yandex.ru/",
			// 		ShortURL:      "short_url",
			// 		CorrelationID: "H3HIkks3",
			// 	},
			// },
		},

		{
			name: "Test GET negative",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				location:    "",
			},
			req: req{
				url: "/N9KHHoG1",
			},
			// store: map[string]*model.URLTb{
			// 	"H3HIkks3": {
			// 		ID:            0,
			// 		OriginalURL:   "https://practicum.yandex.ru/",
			// 		ShortURL:      "short_url",
			// 		CorrelationID: "H3HIkks3",
			// },
			// },
		},
	}

	//
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Request
			request := httptest.NewRequest(http.MethodGet, tt.req.url, nil)
			// Recorder
			response := httptest.NewRecorder()

			// Handler
			h := handler.HandlerURL{CrudSrver: shURL, DelSrver: nil}

			h.ReadURL(response, request)

			// StatusCode
			if response.Code != tt.want.statusCode {
				t.Errorf("%s:\n\texpected: %v\n\tactual: %v", tt.name, tt.want.statusCode, response.Code)
			}

			// Content-Type
			if response.Header().Get("Content-Type") != tt.want.contentType {
				t.Errorf("%s:\n\texpected: %v\n\tactual: %v", tt.name, tt.want.contentType, response.Header().Get("Content-Type"))
			}

			// Location
			if response.Header().Get("Location") != tt.want.location {
				t.Errorf("%s:\n\texpected: %v\n\tactual: %v", tt.name, tt.want.location, response.Header().Get("Location"))
			}
		})
	}
}
