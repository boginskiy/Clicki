package middleware

import "net/http"

// Middleware - .
type Middleware interface {
	WithInfoLogger(http.HandlerFunc) http.HandlerFunc
	WithGzip(http.HandlerFunc) http.HandlerFunc
	Conveyor(http.HandlerFunc) http.HandlerFunc
	WithAuth(http.HandlerFunc) http.HandlerFunc
}
