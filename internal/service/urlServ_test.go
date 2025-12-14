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
	"github.com/boginskiy/Clicki/internal/tester/tfunc"
	"github.com/boginskiy/Clicki/internal/validation"
	"github.com/stretchr/testify/assert"
)

const (
	USERS = 1000
	URL   = "https://practicum.yandex.ru"
)

func BenchmarkRead(b *testing.B) {
	// Logger & Config
	pathToLogg := "test.log"
	logg := logg.NewLogg(pathToLogg, "INFO")

	config := InitConfig()

	// Service
	servShortURL := InitURLServ(logg, config)

	// Request
	UserID := 100
	ctx := context.WithValue(context.Background(), auth.CtxUserID, UserID)
	Method := "GET"
	URL := "/wrs4db6j"

	request := tfunc.PreparRequest(ctx, Method, URL, nil)

	b.ResetTimer() // Обнуление счетчика

	for i := 1; i < b.N; i++ {
		dataByte, _ := servShortURL.Read(request)
		_ = dataByte
	}

	defer tfunc.DeleteTestFiles(pathToLogg)
}

func BenchmarkCreateURL(b *testing.B) {
	// Init
	pathToFile := "test.log"
	logg := logg.NewLogg(pathToFile, "INFO")
	cfg := InitConfig()

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
		request := tfunc.PreparRequest(ctx, Method, URL, body)

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

	defer tfunc.DeleteTestFiles(pathToFile)
}

func TestURLServ(t *testing.T) {
	// Init
	pathToFile := "test.log"
	logg := logg.NewLogg(pathToFile, "INFO")
	cfg := InitConfig()

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
	testTakeUserIDFromCtx(t, ctx, URLServ)
	testEncrypOriginURL(t, URLServ)

	defer tfunc.DeleteTestFiles(pathToFile)
}

func testEncrypOriginURL(t *testing.T, srv *URLServ) {
	res := srv.encrypOriginURL()
	assert.Equal(t, len(res), 8)
}

func testTakeUserIDFromCtx(t *testing.T, ctx context.Context, srv *URLServ) {
	Method := "POST"
	URL := "/"

	// request with userID := 100
	request := tfunc.PreparRequest(ctx, Method, URL, nil)
	assert.Equal(t, srv.takeUserIDFromCtx(request), 100)

	// request without userID
	request2 := tfunc.PreparRequest(context.TODO(), Method, URL, nil)
	assert.NotEqual(t, srv.takeUserIDFromCtx(request2), 100)

}

func testReadSet(t *testing.T, ctx context.Context, srv *URLServ) {
	request := tfunc.PreparRequest(ctx, "GET", "/", nil)
	dataByte, err := srv.ReadSet(request)
	assert.NoError(t, err)
	assert.Greater(t, len(dataByte), 0)
}

func testCreateSet(t *testing.T, ctx context.Context, srv *URLServ) {
	request := tfunc.PreparRequest(ctx, "GET", "/", nil)
	dataByte, err := srv.CreateSet(request)
	assert.NoError(t, err)
	assert.Greater(t, len(dataByte), 0)
}

func testCheckDB(t *testing.T, ctx context.Context, srv *URLServ) {
	Method := "GET"
	URL := "/ping"
	request := tfunc.PreparRequest(ctx, Method, URL, nil)
	dataByte, err := srv.CheckDB(request)
	assert.NoError(t, err)
	assert.Greater(t, len(dataByte), 0)
}

func testRead(t *testing.T, ctx context.Context, srv *URLServ) {
	Method := "GET"
	URL := "/wrs4db6j"

	msg := "check ReadURL in URLServ"
	request := tfunc.PreparRequest(ctx, Method, URL, nil)

	dataByte, err := srv.Read(request)
	if err != nil {
		tfunc.PprintErr(t, msg, err, nil)
	}
	if string(dataByte) != "https://practicum.yandex.ru/" {
		tfunc.PprintErr(t, msg, string(dataByte), "https://practicum.yandex.ru/")
	}
}

func testCreate(t *testing.T, ctx context.Context, srv *URLServ) {
	Body := []byte("https://www.google.com/chrome/")
	Method := "POST"
	URL := "/"

	msg := "check CreateURL in URLServ"
	request := tfunc.PreparRequest(ctx, Method, URL, Body)

	dataByte, err := srv.Create(request)
	if err != nil {
		tfunc.PprintErr(t, msg, err, nil)
	}
	if len(dataByte) == 0 {
		tfunc.PprintErr(t, msg, len(dataByte), "Must be > 0")
	}
}

// InitURLServ.
func InitURLServ(logger logg.Logger, cfg config.Config) *URLServ {
	db, _ := database.NewStoreMap(cfg, logger)
	repo := repository.NewMainRepoMap(cfg, logger, db)

	tfunc.WriteRecord(repo)

	var sub1 = audit.NewFileReceiver(logger, cfg.GetAuditFile(), 1)
	var sub2 = audit.NewServerReceiver(logger, cfg.GetAuditURL(), 2)
	var publisher = audit.NewPublish(sub1, sub2)

	var fancer = preparation.NewFunctions()
	var checker = validation.NewChecker()

	return NewURLServ(cfg, logger, repo, checker, fancer, publisher)
}

// InitConfig.
func InitConfig() *config.Variables {
	return &config.Variables{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080",
		ArgsCLI:       &config.ArgsCLI{},
		ArgsENV: &config.ArgsENV{
			SoftDeleteTime: 10,
			HardDeleteTime: 20,
		},
	}
}
