package app_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/boginskiy/Clicki/internal/auth"
	"github.com/boginskiy/Clicki/internal/logg"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/protocol"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/internal/tester/tfunc"
	"github.com/boginskiy/Clicki/internal/tester/tserv"
)

// Services
var APIURLServ *service.APIURLServ
var URLServ *service.URLServ

// Handler
var APIURLServCreate = APIURLServ.Create
var URLServCreate = URLServ.Create
var URLServRead = URLServ.Read

func ExampleURLServRead() {
	// Init.
	pathToLogg := "test.log"
	logg := logg.NewLogg(pathToLogg, "ERROR")
	config := tserv.InitConfig()

	URLServ := tserv.InitURLServ(logg, config)

	// Request.'
	ctx := context.WithValue(context.Background(), auth.CtxUserID, 100)
	request := tfunc.PreparRequest(ctx, "GET", "/wrs4db6j", nil)
	// Handler.
	body, err := URLServ.Read(request)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(body))
	// Output:
	// https://practicum.yandex.ru/

	defer tfunc.DeleteTestFiles(pathToLogg)
}

func ExampleURLServCreate() {
	pathToLogg := "test.log"

	logg := logg.NewLogg(pathToLogg, "ERROR")
	config := tserv.InitConfig()

	URLServ := tserv.InitURLServ(logg, config)

	// Request.
	request := tfunc.PreparRequest(context.TODO(), "POST", "/", []byte("https://practicum.yandex.ru"))
	// Handler.
	body, err := URLServ.Create(request)
	if err != nil {
		log.Fatalln(err)
	}

	result := string(body)[:len(string(body))-9]

	fmt.Println(result)
	// Output:
	// http://localhost:8080

	defer tfunc.DeleteTestFiles(pathToLogg)
}

func ExampleAPIURLServCreate() {
	// Init.
	pathToLogg := "test.log"
	logg := logg.NewLogg(pathToLogg, "ERROR")
	config := tserv.InitConfig()

	TestAPIURLServ := tserv.InitAPIURLServ(logg, config)

	// Request.
	ctx := context.WithValue(context.Background(), auth.CtxUserID, 100)
	body := []byte(`{"url": "https://leetcode.com/"}`)
	request := tfunc.PreparRequest(ctx, "POST", "/shorten", body)

	// Protocol
	funcer := prep.NewFunctions(logg)
	protHTTP := protocol.NewProtocolHTTP(funcer)

	// Handler.
	body, err := TestAPIURLServ.Create(request.Context(), protHTTP, request)
	if err != nil {
		log.Fatalln(err)
	}

	mapp := map[string]any{}
	json.Unmarshal(body, &mapp)

	url := mapp["result"].(string)

	fmt.Println(url[:len(url)-9])
	// Output:
	// http://localhost:8080

	defer tfunc.DeleteTestFiles(pathToLogg)
}
