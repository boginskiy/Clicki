package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/boginskiy/Clicki/internal/auth"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/internal/tester"
)

const (
	USERS = 1000
	URL   = "https://practicum.yandex.ru"
)

func BenchmarkRead(b *testing.B) {
	// Logger & Config
	pathToLogg := "test.log"
	logg := logg.NewLogg(pathToLogg, "INFO")

	config := tester.InitConfig()

	// Service
	servShortURL := tester.InitURLServ(logg, config)

	// Request
	UserID := 100
	ctx := context.WithValue(context.Background(), auth.CtxUserID, UserID)
	Method := "GET"
	URL := "/wrs4db6j"

	request := tester.PreparRequest(ctx, Method, URL, nil)

	b.ResetTimer() // Обнуление счетчика

	for i := 1; i < b.N; i++ {
		dataByte, _ := servShortURL.Read(request)
		_ = dataByte
	}

	defer tester.DeleteTestFiles(pathToLogg)
}

func BenchmarkCreateURL(b *testing.B) {
	// Init
	pathToFile := "test.log"
	logg := logg.NewLogg(pathToFile, "INFO")
	cfg := tester.InitConfig()

	// Service
	URLServ := tester.InitURLServ(logg, cfg)

	// Request
	body := []byte(fmt.Sprintf("%s%d%s", "https://practicum.yandex-", 0, ".ru"))
	Method := "POST"
	UserID := 100
	URL := "/"

	b.ResetTimer()    // Обнуление счетчика
	start, N := 0, 20 // Счетчик | Количество URL на пользователя

	for i := 1; i < b.N; i++ {
		ctx := context.WithValue(context.Background(), auth.CtxUserID, UserID)
		request := tester.PreparRequest(ctx, Method, URL, body)

		dataByte, _ := URLServ.Create(request)
		_ = dataByte

		// останавливаем таймер
		b.StopTimer()

		// После N добавленных URL меняем пользователя
		if start > N {
			start = 0
			UserID++
		} else {
			start++
		}

		// Body меняем каждую итерацию, чтобы максимально полно оценить работу метода
		body = []byte(fmt.Sprintf("%s%d%s", "https://practicum.yandex-", i, ".ru"))
		// возобновляем таймер
		b.StartTimer()
	}

	defer tester.DeleteTestFiles(pathToFile)
}

func TestURLServ(t *testing.T) {

	// Init
	pathToFile := "test.log"
	logg := logg.NewLogg(pathToFile, "INFO")
	cfg := tester.InitConfig()

	// Service
	URLServ := tester.InitURLServ(logg, cfg)

	// Ctx
	userID := 100
	ctx := context.WithValue(context.Background(), auth.CtxUserID, userID)

	testCreate(t, ctx, URLServ)
	testRead(t, ctx, URLServ)

	defer tester.DeleteTestFiles(pathToFile)
}

func testRead(t *testing.T, ctx context.Context, srv *service.URLServ) {
	Method := "GET"
	URL := "/wrs4db6j"

	msg := "check ReadURL in service.URLServ"
	request := tester.PreparRequest(ctx, Method, URL, nil)

	dataByte, err := srv.Read(request)
	if err != nil {
		tester.PprintErr(t, msg, err, nil)
	}
	if string(dataByte) != "https://practicum.yandex.ru/" {
		tester.PprintErr(t, msg, string(dataByte), "https://practicum.yandex.ru/")
	}
}

func testCreate(t *testing.T, ctx context.Context, srv *service.URLServ) {
	Body := []byte("https://www.google.com/chrome/")
	Method := "POST"
	URL := "/"

	msg := "check CreateURL in service.URLServ"
	request := tester.PreparRequest(ctx, Method, URL, Body)

	dataByte, err := srv.Create(request)
	if err != nil {
		tester.PprintErr(t, msg, err, nil)
	}
	if len(dataByte) == 0 {
		tester.PprintErr(t, msg, len(dataByte), "Must be > 0")
	}
}
