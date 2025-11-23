package tester

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/boginskiy/Clicki/internal/service"
)

func PprintErr(t *testing.T, msg string, oneArg, twoArg any) {
	t.Errorf("%s: reality: %v, expect: %v\n", msg, oneArg, twoArg)
}

func PreparRequest(ctx context.Context, method, url string, body []byte) *http.Request {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Prepar Request: %v != %v\n", err, nil)
	}
	return req
}

func CheckMethodOfServ(t *testing.T, msg string, data []byte, err error) {
	if err != nil {
		PprintErr(t, msg, err, nil)
	}
	if len(data) == 0 {
		PprintErr(t, msg, len(data), service.LONG)
	}
}
