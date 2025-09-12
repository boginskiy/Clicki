package service_test

import (
	"testing"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/logger"
	"github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/internal/validation"
)

var infoLog = logger.NewLogg("Test.log", "INFO")
var kwargs = &config.Variables{
	ServerAddress: "localhost:8080",
	BaseURL:       "http://localhost:8081",
}

var repo, _ = db.NewStoreMap(kwargs, infoLog)

var extraFuncer = preparation.NewExtraFunc()
var checker = validation.NewChecker()

var ShURL = service.NewShortURL(repo, infoLog, checker, extraFuncer)

func TestEncryptionLongURL(t *testing.T) {
	name := "Check EncryptionLongURL from ProURL"
	shortURL := "dcJd743D"

	// Проверка длины
	expected := service.LONG
	if expected != len(shortURL) {
		t.Errorf("Test 1 >> %s > expected %v actual > %v", name, expected, len(shortURL))
	}

	// Проверка по регулярному выражению
	expected2 := true
	actual2 := ShURL.Checker.CheckUpPath(shortURL)
	if ShURL.Checker.CheckUpPath("/"+shortURL) != true {
		t.Errorf("Test 2 >> %s > expected %v actual > %v", name, expected2, actual2)
	}
}
