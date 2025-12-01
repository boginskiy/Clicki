package middleware

import "net/http"

// ExReqWriter - extra http.Request.
type ExReqWriter struct {
	*http.Request
	UserID int
}
