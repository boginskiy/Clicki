package register

import "net/http"

type Register interface {
	Registration(*http.Request, *int) ([]byte, *http.Cookie, error)
	Authentication(*http.Request, *int) ([]byte, *http.Cookie, error)
}
