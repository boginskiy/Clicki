package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	endpoint := "http://localhost:8080/"

	// Контейнер данных для запроса
	data := url.Values{}

	// Приглашение в консоли
	fmt.Println("Введите длинный URL")
	// Открываем потоковое чтение из консоли
	reader := bufio.NewReader(os.Stdin)

	// Читаем строку из консоли
	long, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	long = strings.TrimSuffix(long, "\n")
	// Заполянем контейнер данными
	data.Set("url", long)

	// Добавляем HTTP-клиента
	client := &http.Client{}

	// пишем запрос
	// запрос методом POST должен, помимо заголовков, содержать тело
	// тело должно быть источником потокового чтения io.Reader

	request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	// в заголовках запроса указываем кодировку
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// отправляем запрос и получаем ответ
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	// Выводим код ответа
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()

	// читаем поток из тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// И печатаем его
	fmt.Println(string(body))
}
