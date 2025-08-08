package service_test

import (
	"testing"

	"github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/service"
)

var database = db.NewDbStore()
var ShortingURL = service.NewShorteningURL(database)

func TestEncryptionLongURL(t *testing.T) {
	name := "Check EncryptionLongURL from ProURL"
	imitationPath := ShortingURL.EncryptionLongURL()

	// Проверка длины
	expected := service.LONG
	if expected != len(imitationPath) {
		t.Errorf("Test 1 >> %s > expected %v actual > %v", name, expected, len(imitationPath))
	}

	// Проверка по регулярному выражению
	expected2 := true
	actual2 := ShortingURL.CheckUpPath(imitationPath)
	if ShortingURL.CheckUpPath("/"+imitationPath) != true {
		t.Errorf("Test 2 >> %s > expected %v actual > %v", name, expected2, actual2)
	}
}
