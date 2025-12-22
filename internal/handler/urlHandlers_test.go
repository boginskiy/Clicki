package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/boginskiy/Clicki/internal/auth"
	"github.com/boginskiy/Clicki/internal/handler"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/model"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/internal/tester/tfunc"
	"github.com/boginskiy/Clicki/internal/tester/tserv"
)

func TestHandlerURL(t *testing.T) {
	// Logger & Config.
	pathToLogg := "test.log"
	logg := logg.NewLogg(pathToLogg, "INFO")

	config := tserv.InitConfig()

	URLServ := tserv.InitURLServ(logg, config)

	// Testing
	testCreateRecord(t, URLServ)
	testReadRecord(t, URLServ)

	defer tfunc.DeleteTestFiles(pathToLogg)
}

func testCreateRecord(t *testing.T, URLServ *service.URLServ) {
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
			name: "test positive create record with URLHandlers",
			want: want{
				contentType: "text/plain",
				statusCode:  201,
			},
			req: req{
				url:  "/",
				body: "https://www.google.com/chrome/",
				host: "localhost:8080",
			},
		},
		{
			name: "test negative create record with URLHandlers",
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Request
			ctx := context.WithValue(context.Background(), auth.CtxUserID, 100)
			body := []byte(tt.req.body)
			request := tfunc.PreparRequest(ctx, http.MethodPost, tt.req.url, body)

			// Recorder
			response := httptest.NewRecorder()

			// URLHandlers
			URLHandlers := handler.NewURLHandlers(URLServ)
			URLHandlers.Create(response, request)

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
		})
	}
}

func testReadRecord(t *testing.T, URLServ *service.URLServ) {
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
			name: "test positive read record with URLHandlers",
			want: want{
				contentType: "",
				statusCode:  307,
				location:    "https://practicum.yandex.ru/",
			},
			req: req{
				url: "/wrs4db6j",
			},
		},
		{
			name: "test negative read record with URLHandlers",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
				location:    "",
			},
			req: req{
				url: "/N9KHHoG1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Request
			ctx := context.WithValue(context.Background(), auth.CtxUserID, 100)
			request := tfunc.PreparRequest(ctx, http.MethodGet, tt.req.url, nil)

			// Recorder
			response := httptest.NewRecorder()

			// URLHandlers
			URLHandlers := handler.NewURLHandlers(URLServ)
			URLHandlers.Read(response, request)

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
