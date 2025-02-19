package main

import (
	"net/http"
	"os"
	"testing"
	"time"
)

func TestRunServer(t *testing.T) {
	// Устанавливаем переменную окружения ADDRESS
	os.Setenv("ADDRESS", "localhost:8080")

	// Запускаем сервер в отдельной горутине
	go func() {
		main() // Это ваша функция, которая запускает сервер
	}()

	// Даем серверу время на запуск
	time.Sleep(5 * time.Second)

	// Отправляем GET-запрос на сервер
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Fatalf("Ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	os.Unsetenv("ADDRESS")

}
