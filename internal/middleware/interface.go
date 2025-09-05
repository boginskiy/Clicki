package middleware

import "net/http"

type Middlewarer interface {
	WithInfoLogger(http.HandlerFunc) http.HandlerFunc
	WithGzip(http.HandlerFunc) http.HandlerFunc
	Conveyor(http.HandlerFunc) http.HandlerFunc
}
