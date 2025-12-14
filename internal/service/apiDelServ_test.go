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
	"github.com/boginskiy/Clicki/internal/tester"
	"github.com/boginskiy/Clicki/internal/tester/testinit"
	"github.com/boginskiy/Clicki/internal/validation"
	"github.com/stretchr/testify/assert"
)

func TestApiDelServ(t *testing.T) {
	// Init.
	pathToLogg := "test.log"
	// Cong & Logg
	config, logg := testinit.InitConfAndLogg(pathToLogg)

	// DB, Repo
	db, _ := database.NewStoreMap(config, logg)
	repo := repository.NewMainRepoMap(config, logg, db)

	tester.WriteRecord(repo)

	// Serv
	APIDelServ := NewAPIDelServ(context.TODO(), config, logg, repo)
	defer tester.DeleteTestFiles(pathToLogg)

	// Testing
	ctx, cancel := context.WithCancel(context.Background())
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

}

func testDeleteSet(wg *sync.WaitGroup, t *testing.T, serv *APIDelServ) {
	defer wg.Done()
	// Request
	ctx := context.WithValue(context.Background(), auth.CtxUserID, 100)

	dataBody := []string{"wrs4db6j"}
	request := tester.PreparRequest(ctx, "DELETE", "/user/urls", tester.Serialization(dataBody))

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

	tester.WriteRecord(repo)

	var sub1 = audit.NewFileReceiver(logger, cfg.GetAuditFile(), 1)
	var sub2 = audit.NewServerReceiver(logger, cfg.GetAuditURL(), 2)
	var publisher = audit.NewPublish(sub1, sub2)

	var fancer = preparation.NewFunctions()
	var checker = validation.NewChecker()

	return NewAPIURLServ(cfg, logger, repo, checker, fancer, publisher)
}
