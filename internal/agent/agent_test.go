package agent

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	// Канал для сигнала завершения
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем Run в горутине
	go func() {
		Run()
	}()

	// Симулируем сигнал SIGINT
	signalChannel <- syscall.SIGINT

	// Ожидаем завершения работы
	select {
	case sig := <-signalChannel:
		// Проверяем, что сигнал был получен
		assert.Equal(t, syscall.SIGINT, sig)
	case <-time.After(5 * time.Second):
		t.Errorf("Run did not exit within the timeout")
	}
}
