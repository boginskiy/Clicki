package middleware

import "net/http"

type ExReqWriter struct {
	*http.Request
	UserID int
}

type contextKey struct{}

var CtxUserID = contextKey{}
