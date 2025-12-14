package service

import (
	"context"
	"sync"
	"testing"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/audit"
	"github.com/boginskiy/Clicki/internal/auth"
	"github.com/boginskiy/Clicki/internal/database"
	"github.com/boginskiy/Clicki/internal/layers"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/internal/repository/utils"
	"github.com/boginskiy/Clicki/internal/tester/tfunc"
	"github.com/boginskiy/Clicki/internal/validation"
	"github.com/stretchr/testify/assert"
)

func TestApiDelServ(t *testing.T) {
	// Init
	pathToFile := "test.log"
	logg := logg.NewLogg(pathToFile, "INFO")
	config := InitConfig()

	// DB, Repo
	db, _ := database.NewStoreMap(config, logg)
	repo := repository.NewMainRepoMap(config, logg, db)

	tfunc.WriteRecord(repo)

	// Serv
	APIDelServ := NewAPIDelServ(context.TODO(), config, logg, repo)
	defer tfunc.DeleteTestFiles(pathToFile)

	// Testing
	userID := 100
	ctx := context.WithValue(context.Background(), auth.CtxUserID, userID)
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup

	wg.Add(2)
	go testDeleteSet(&wg, t, APIDelServ) // Result: d.delMessChan <- *delMessage
	go testStepByStepDelMessages(ctx, &wg, t, APIDelServ)
	// time.Sleep(50 * time.Millisecond)

	/*
		Не совсем корректный тест, потому что логика завязанная на канале не успевает отработать и мы сразу
		делаем cancel.

		Для репозитория типа Map по другому быть не может так как этот репозиторий не реализует логику для обработки
		канальной логики. Будет ошибка.

		Потенциально это место может работать непредсказуемо.

	*/
	cancel()
	wg.Wait()

	testSendHardDeletion(t, APIDelServ)
	testTakeUserIDFromCtx2(t, ctx, APIDelServ)
	testSendSoftDeletion(t, ctx, APIDelServ)

}

/*
Тестирую только так, чтобы len(data) == 0, потому что для DB == Map
d.Repo.MarkRecords == nil
*/
func testSendSoftDeletion(t *testing.T, ctx context.Context, srv *APIDelServ) {
	// dmess := utils.DelMessage{
	// 	ListCorrelID: []string{"wrs4db6j"},
	// 	UserID:       100}

	data := []utils.DelMessage{}
	isDel := true

	res := srv.sendSoftDeletion(data, &isDel)
	assert.Equal(t, len(res), 0)
}

func testTakeUserIDFromCtx2(t *testing.T, ctx context.Context, srv *APIDelServ) {
	Method := "POST"
	URL := "/"

	// request with userID := 100
	request := tfunc.PreparRequest(ctx, Method, URL, nil)
	assert.Equal(t, srv.takeUserIDFromCtx(request), 100)

	// request without userID
	request2 := tfunc.PreparRequest(context.TODO(), Method, URL, nil)
	assert.NotEqual(t, srv.takeUserIDFromCtx(request2), 100)

}

func testSendHardDeletion(t *testing.T, serv *APIDelServ) {
	var isTrue bool = true
	assert.Equal(t, serv.sendHardDeletion(&isTrue), false)
}

func testDeleteSet(wg *sync.WaitGroup, t *testing.T, serv *APIDelServ) {
	defer wg.Done()
	// Request
	ctx := context.WithValue(context.Background(), auth.CtxUserID, 100)

	dataBody := []string{"wrs4db6j"}
	request := tfunc.PreparRequest(ctx, "DELETE", "/user/urls", tfunc.Serialization(dataBody))

	_, err := serv.DeleteSet(request)
	assert.NoError(t, err)

}

func testStepByStepDelMessages(ctx context.Context, wg *sync.WaitGroup, t *testing.T, serv *APIDelServ) {
	defer wg.Done()
	serv.StepByStepDelMessages(ctx)
}

// InitAPIURLServ
func InitAPIURLServ(logger logg.Logger, cfg config.Config) *APIURLServ {
	layers := layers.NewLayers(cfg, logger)
	db := layers.NewLayerDB()
	repo := layers.NewLayerRepo(db)

	tfunc.WriteRecord(repo)

	var sub1 = audit.NewFileReceiver(logger, cfg.GetAuditFile(), 1)
	var sub2 = audit.NewServerReceiver(logger, cfg.GetAuditURL(), 2)
	var publisher = audit.NewPublish(sub1, sub2)

	var fancer = preparation.NewFunctions()
	var checker = validation.NewChecker()

	return NewAPIURLServ(cfg, logger, repo, checker, fancer, publisher)
}
