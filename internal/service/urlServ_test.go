package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/audit"
	"github.com/boginskiy/Clicki/internal/auth"
	"github.com/boginskiy/Clicki/internal/database"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/internal/tester"
	"github.com/boginskiy/Clicki/internal/validation"

	"github.com/boginskiy/Clicki/internal/tester/testinit"
	"github.com/stretchr/testify/assert"
)

// Проверить ,что нет цикленности когда мы не юзаем в уровне tester сущности
// Надо тестировать сервисы далее.

const (
	USERS = 1000
	URL   = "https://practicum.yandex.ru"
)

func BenchmarkRead(b *testing.B) {
	// Logger & Config
	pathToLogg := "test.log"
	logg := logg.NewLogg(pathToLogg, "INFO")

	config := testinit.InitConfig()

	// Service
	servShortURL := InitURLServ(logg, config)

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
	cfg := testinit.InitConfig()

	// Service
	URLServ := InitURLServ(logg, cfg)

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
	cfg := testinit.InitConfig()

	// Service
	URLServ := InitURLServ(logg, cfg)

	// Ctx
	userID := 100
	ctx := context.WithValue(context.Background(), auth.CtxUserID, userID)

	// Public method
	testCreateSet(t, ctx, URLServ)
	testReadSet(t, ctx, URLServ)
	testCheckDB(t, ctx, URLServ)
	testCreate(t, ctx, URLServ)
	testRead(t, ctx, URLServ)

	// Private method
	testtakeUserIDFromCtx(t, ctx, URLServ)

	defer tester.DeleteTestFiles(pathToFile)
}

// func (s *URLServ) takeUserIDFromCtx(req *http.Request) int {
// 	UserID, ok := req.Context().Value(auth.CtxUserID).(int)
// 	if !ok || UserID <= 0 {
// 		s.Logg.RaiseError(ErrUserIDNotValid, "URLServ.takeUserIDFromCtx>CtxUserID", nil)
// 	}
// 	return UserID
// }

func testtakeUserIDFromCtx(t *testing.T, ctx context.Context, srv *URLServ) {

}

func testReadSet(t *testing.T, ctx context.Context, srv *URLServ) {
	request := tester.PreparRequest(ctx, "GET", "/", nil)
	dataByte, err := srv.ReadSet(request)
	assert.NoError(t, err)
	assert.Greater(t, len(dataByte), 0)
}

func testCreateSet(t *testing.T, ctx context.Context, srv *URLServ) {
	request := tester.PreparRequest(ctx, "GET", "/", nil)
	dataByte, err := srv.CreateSet(request)
	assert.NoError(t, err)
	assert.Greater(t, len(dataByte), 0)
}

func testCheckDB(t *testing.T, ctx context.Context, srv *URLServ) {
	Method := "GET"
	URL := "/ping"
	request := tester.PreparRequest(ctx, Method, URL, nil)
	dataByte, err := srv.CheckDB(request)
	assert.NoError(t, err)
	assert.Greater(t, len(dataByte), 0)
}

func testRead(t *testing.T, ctx context.Context, srv *URLServ) {
	Method := "GET"
	URL := "/wrs4db6j"

	msg := "check ReadURL in URLServ"
	request := tester.PreparRequest(ctx, Method, URL, nil)

	dataByte, err := srv.Read(request)
	if err != nil {
		tester.PprintErr(t, msg, err, nil)
	}
	if string(dataByte) != "https://practicum.yandex.ru/" {
		tester.PprintErr(t, msg, string(dataByte), "https://practicum.yandex.ru/")
	}
}

func testCreate(t *testing.T, ctx context.Context, srv *URLServ) {
	Body := []byte("https://www.google.com/chrome/")
	Method := "POST"
	URL := "/"

	msg := "check CreateURL in URLServ"
	request := tester.PreparRequest(ctx, Method, URL, Body)

	dataByte, err := srv.Create(request)
	if err != nil {
		tester.PprintErr(t, msg, err, nil)
	}
	if len(dataByte) == 0 {
		tester.PprintErr(t, msg, len(dataByte), "Must be > 0")
	}
}

// InitURLServ init serv.
func InitURLServ(logger logg.Logger, cfg config.Config) *URLServ {
	db, _ := database.NewStoreMap(cfg, logger)
	repo := repository.NewMainRepoMap(cfg, logger, db)

	tester.WriteRecord(repo)

	var sub1 = audit.NewFileReceiver(logger, cfg.GetAuditFile(), 1)
	var sub2 = audit.NewServerReceiver(logger, cfg.GetAuditURL(), 2)
	var publisher = audit.NewPublish(sub1, sub2)

	var fancer = preparation.NewFunctions()
	var checker = validation.NewChecker()

	return NewURLServ(cfg, logger, repo, checker, fancer, publisher)
}
