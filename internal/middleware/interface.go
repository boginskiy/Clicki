package middleware

import "net/http"

type Middlewarer interface {
	WithInfoLogger(http.HandlerFunc) http.HandlerFunc
}
