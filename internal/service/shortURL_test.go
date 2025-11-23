package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/audit"
	"github.com/boginskiy/Clicki/internal/auther"
	"github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/model"
	"github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/internal/tester"
	"github.com/boginskiy/Clicki/internal/validation"
)

const (
	USERS = 1000
	URL   = "https://practicum.yandex.ru"
)

func BenchmarkReadURL(b *testing.B) {
	// Init
	logger := InitLogg("test.log", "INFO")
	kwargs := InitKwargs()
	db := InitDB()

	// Service
	servShortURL := Init(logger, kwargs, db)

	// Request
	Method := "GET"
	UserID := 100
	URL := "/H3HIkks3"

	// Ctx
	ctx := context.WithValue(context.Background(), auther.CtxUserID, UserID)
	request := tester.PreparRequest(ctx, Method, URL, nil)

	b.ResetTimer() // Обнуление счетчика

	for i := 1; i < b.N; i++ {
		dataByte, _ := servShortURL.ReadURL(request)
		_ = dataByte
	}
}

func BenchmarkCreateURL(b *testing.B) {
	// Init
	logger := InitLogg("test.log", "INFO")
	kwargs := InitKwargs()
	db := InitDB()

	// Service
	servShortURL := Init(logger, kwargs, db)

	// Request
	body := []byte(fmt.Sprintf("%s%d%s", "https://practicum.yandex-", 0, ".ru"))
	Method := "POST"
	UserID := 100
	URL := "/"

	b.ResetTimer()    // Обнуление счетчика
	start, N := 0, 20 // Счетчик | Количество URL на пользователя

	for i := 1; i < b.N; i++ {
		ctx := context.WithValue(context.Background(), auther.CtxUserID, UserID)
		request := tester.PreparRequest(ctx, Method, URL, body)

		dataByte, _ := servShortURL.CreateURL(request)
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
}

func TestShortURL(t *testing.T) {
	// Init
	logger := InitLogg("test.log", "INFO")
	kwargs := InitKwargs()
	db := InitDB()

	// Service
	servShortURL := Init(logger, kwargs, db)

	// Ctx
	userID := 100
	ctx := context.WithValue(context.Background(), auther.CtxUserID, userID)

	testCreateURL(t, ctx, servShortURL)
	testReadURL(t, ctx, servShortURL)
}

func testReadURL(t *testing.T, ctx context.Context, srv *service.ShortURL) {
	Method := "GET"
	URL := "/H3HIkks3"

	msg := "check ReadURL in service.ShortURL"
	request := tester.PreparRequest(ctx, Method, URL, nil)

	dataByte, err := srv.ReadURL(request)
	if err != nil {
		tester.PprintErr(t, msg, err, nil)
	}
	if string(dataByte) != "https://giga.chat/" {
		tester.PprintErr(t, msg, string(dataByte), "https://giga.chat/")
	}
}

func testCreateURL(t *testing.T, ctx context.Context, srv *service.ShortURL) {
	Body := []byte("https://practicum.yandex.ru")
	Method := "POST"
	URL := "/"

	msg := "check CreateURL in service.ShortURL"
	request := tester.PreparRequest(ctx, Method, URL, Body)

	dataByte, err := srv.CreateURL(request)
	if err != nil {
		tester.PprintErr(t, msg, err, nil)
	}
	if len(dataByte) == 0 {
		tester.PprintErr(t, msg, len(dataByte), "Must be > 0")
	}
}

func Init(logger logg.Logger, kwargs config.VarGetter, db db.DBer) *service.ShortURL {
	var repo, _ = repository.NewRepositoryMapURL(kwargs, db)

	var sub1 = audit.NewFileReceiver(logger, kwargs.GetAuditFile(), 1)
	var sub2 = audit.NewServerReceiver(logger, kwargs.GetAuditURL(), 2)
	var publisher = audit.NewPublish(sub1, sub2)

	var core = service.NewCoreService(kwargs, logger, repo, publisher)
	var extraFuncer = preparation.NewExtraFunc()
	var checker = validation.NewChecker()

	return service.NewShortURL(core, repo, checker, extraFuncer)
}

func InitKwargs() config.VarGetter {
	return &config.Variables{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8081",
	}
}

func InitLogg(name, mod string) logg.Logger {
	return logg.NewLogg(name, mod)
}

func InitDB() db.DBer {
	// Map
	row := model.NewURLTb(1, "H3HIkks3", "https://giga.chat/", "http://localhost:8081/H3HIkks3", 100)
	return &db.StoreMap{Store: map[string]*model.URLTb{"H3HIkks3": row}}
	// Файл
	// БД
}
