package server_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/cmd/server"
	"github.com/boginskiy/Clicki/internal/auth"
	"github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/logg"
	"github.com/boginskiy/Clicki/internal/model"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/internal/tester"
	"github.com/boginskiy/Clicki/internal/validation"
)

// Services
var APIShortURL *service.APIShortURL
var ShortURL *service.ShortURL

// Handler
var APIShortURLCreateURL = APIShortURL.CreateURL
var ShortURLCreateURL = ShortURL.CreateURL
var ShortURLReadURL = ShortURL.ReadURL

func ExampleShortURLReadURL() {
	// Init.
	TShortURL := InitServices1()

	// Request.'
	ctx := context.WithValue(context.Background(), auth.CtxUserID, 100)
	request := tester.PreparRequest(ctx, "GET", "/H3HIkks3", nil)
	// Handler.
	body, err := TShortURL.ReadURL(request)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(body))
	// Output:
	// https://practicum.yandex.ru/
}

func ExampleShortURLCreateURL() {
	// Init.
	TshortURL := InitServices1()
	// Request.
	request := tester.PreparRequest(context.TODO(), "POST", "/", []byte("https://practicum.yandex.ru"))
	// Handler.
	body, err := TshortURL.CreateURL(request)
	if err != nil {
		log.Fatalln(err)
	}

	result := string(body)[:len(string(body))-9]

	fmt.Println(result)
	// Output:
	// http://localhost:8080
}

func ExampleAPIShortURLCreateURL() {
	// Init.
	TapiShortURL := InitServices2()
	// Request.
	request := tester.PreparRequest(context.TODO(), "POST", "/shorten", []byte(`{"url": "https://leetcode.com/"}`))
	// Handler.
	body, err := TapiShortURL.CreateURL(request)
	if err != nil {
		log.Fatalln(err)
	}

	mapp := map[string]any{}
	json.Unmarshal(body, &mapp)

	url := mapp["result"].(string)

	fmt.Println(url[:len(url)-9])
	// Output:
	// http://localhost:8080
}

func InitServices2() *service.APIShortURL {
	// Some part.
	logg := logg.NewLogg("test.log", "ERROR")
	// config := config.NewVariables(logg)
	config := InitConfig()
	checker := validation.NewChecker()
	exFunc := prep.NewExtraFunc()

	// Db.
	layers := server.NewLayers(config, logg)
	db := layers.NewLayerDB()

	// add test data
	db = UpdateDB(db)

	repo := layers.NewLayerRepo(db)

	// Service.
	CoreServ := service.NewCoreService(config, logg, repo, nil)
	return service.NewAPIShortURL(CoreServ, repo, checker, exFunc)
}

func InitServices1() *service.ShortURL {
	// Some part.
	logg := logg.NewLogg("test.log", "ERROR")
	// config := config.NewVariables(logg)
	config := InitConfig()
	checker := validation.NewChecker()
	exFunc := prep.NewExtraFunc()

	// Db.
	layers := server.NewLayers(config, logg)
	db := layers.NewLayerDB()

	// add test data
	db = UpdateDB(db)

	repo := layers.NewLayerRepo(db)

	// Service.
	CoreServ := service.NewCoreService(config, logg, repo, nil)
	return service.NewShortURL(CoreServ, repo, checker, exFunc)
}

func UpdateDB(db db.DataBase) db.DataBase {
	switch v := db.GetDB().(type) {
	case map[string]*model.URLTb:

		record := &model.URLTb{
			ID:            0,
			OriginalURL:   "https://practicum.yandex.ru/",
			ShortURL:      "short_url",
			CorrelationID: "H3HIkks3",
			UserID:        100,
		}

		v["H3HIkks3"] = record
	}
	return db
}

func InitConfig() config.Config {
	return &config.Conf{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080",
	}
}
