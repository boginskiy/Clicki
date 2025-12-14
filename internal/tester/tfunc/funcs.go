package tfunc

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/boginskiy/Clicki/internal/model"
	"github.com/boginskiy/Clicki/internal/repository"
	"golang.org/x/tools/go/analysis"
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

// DeleteTestFiles is Exta function.
func DeleteTestFiles(paths ...string) {
	for _, path := range paths {
		err := os.Remove(path)
		if err != nil {
			log.Fatalf("Error when deleting: %v", err)
		}
	}
}

// WriteRecord is Exta function for add record for tests.
func WriteRecord(repo repository.Repository) {
	record := model.NewURLTb(
		1,
		"wrs4db6j",
		"https://practicum.yandex.ru/",
		"http://localhost:8080/wrs4db6j",
		100)

	repo.CreateRecord(context.TODO(), record)
}

// CreateAnalyzer
func CreateAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "Test",
		Doc:  "check for test",
		Run:  nil,
	}
}

func Serialization(obj any) []byte {
	res, err := json.Marshal(obj)
	if err != nil {
		log.Fatal(err)
	}
	return res
}
