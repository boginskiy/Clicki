package service_test

import (
	"testing"

	"github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/db2"
	"github.com/boginskiy/Clicki/internal/logger"
	"github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/internal/validation"
)

var infoLog = logger.NewLogg("Test.log", "INFO")
var dbase2 = &db2.ConnDB{}

var fileWorker, _ = db.NewFileWorking("test")
var dbase = db.NewDBStore(fileWorker)

var extraFuncer = preparation.NewExtraFunc()
var checker = validation.NewChecker()

var ShURL = service.NewShortURL(dbase, dbase2, infoLog, checker, extraFuncer)

func TestEncryptionLongURL(t *testing.T) {
	name := "Check EncryptionLongURL from ProURL"
	imitationPath := "dcJd743D"

	// Проверка длины
	expected := service.LONG
	if expected != len(imitationPath) {
		t.Errorf("Test 1 >> %s > expected %v actual > %v", name, expected, len(imitationPath))
	}

	// Проверка по регулярному выражению
	expected2 := true
	actual2 := ShURL.Checker.CheckUpPath(imitationPath)
	if ShURL.Checker.CheckUpPath("/"+imitationPath) != true {
		t.Errorf("Test 2 >> %s > expected %v actual > %v", name, expected2, actual2)
	}
}
