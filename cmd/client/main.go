package main

import (
	"bufio"
	"fmt"
	"io"
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
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
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
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// в заголовках запроса указываем кодировку
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// отправляем запрос и получаем ответ
	response, err := client.Do(request)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Выводим код ответа
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()

	// читаем поток из тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// И печатаем его
	fmt.Println(string(body))
}
