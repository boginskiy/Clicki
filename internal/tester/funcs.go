package tester

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"testing"
)

// PprintErr - prity print about errors.
func PprintErr(t *testing.T, msg string, oneArg, twoArg any) {
	t.Errorf("%s: reality: %v, expect: %v\n", msg, oneArg, twoArg)
}

// PreparRequest - .
func PreparRequest(ctx context.Context, method, url string, body []byte) *http.Request {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Prepar Request: %v != %v\n", err, nil)
	}
	return req
}
