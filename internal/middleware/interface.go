package middleware

import "net/http"

type Trustee interface {
	WithTrustedSubnet(http.HandlerFunc) http.HandlerFunc
}

// Middleware - .
type Middleware interface {
	Trustee
	WithLogg(http.HandlerFunc) http.HandlerFunc
	WithGzip(http.HandlerFunc) http.HandlerFunc
	Conveyor(http.HandlerFunc) http.HandlerFunc
	WithAuth(http.HandlerFunc) http.HandlerFunc
}
