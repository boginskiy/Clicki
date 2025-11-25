package server_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/internal/tester"
)

var APIShortURL = service.APIShortURL{}
var ShortURL = service.ShortURL{}

var APIShortURLCreateURL = APIShortURL.CreateURL
var ShortURLCreateURL = ShortURL.CreateURL

func ExampleAPIShortURLCreateURL() {
	// Init.
	TapiShortURL := tester.InitApiShortURL()
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

func ExampleShortURLCreateURL() {
	// Init.
	TshortURL := tester.InitShortURL()
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
