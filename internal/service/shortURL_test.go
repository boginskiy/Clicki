package service_test

import (
	"testing"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/audit"
	"github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/internal/validation"
)

var infoLog = logg.NewLogg("Test.log", "INFO")
var kwargs = &config.Variables{
	ServerAddress: "localhost:8080",
	BaseURL:       "http://localhost:8081",
}

var dber, _ = db.NewStoreMap(kwargs, infoLog)
var repo, _ = repository.NewRepositoryMapURL(kwargs, dber)

// Publisher
var sub1 = audit.NewFileReceiver(infoLog, kwargs.GetAuditFile(), 1)
var sub2 = audit.NewServerReceiver(infoLog, kwargs.GetAuditURL(), 2)
var publisher = audit.NewPublish(sub1, sub2)

var core = service.NewCoreService(kwargs, infoLog, repo, publisher)
var extraFuncer = preparation.NewExtraFunc()
var checker = validation.NewChecker()

var ShURL = service.NewShortURL(core, repo, checker, extraFuncer)

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
