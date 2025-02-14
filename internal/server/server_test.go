package server_test

import (
	"go-metrics-alerting/internal/server"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	// Запуск функции Run в горутине
	go func() {
		server.Run()
	}()

	// Ожидание 2 секунды (для симуляции работы сервера)
	time.Sleep(2 * time.Second)

	// Завершаем тест, ничего не проверяя
	t.Log("Test Run completed")
}
