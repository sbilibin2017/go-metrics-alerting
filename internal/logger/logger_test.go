package logger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func InitLogger(logLevel string) {
	if logLevel == "" {
		logLevel = "INFO"
	}

	// Конфигурируем логгер для JSON-формата
	zapConfig := zap.NewProductionConfig()

	// Устанавливаем уровень логирования в зависимости от конфигурации
	switch logLevel {
	case "DEBUG":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "ERROR":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Создаем логгер с указанной конфигурацией
	var err error
	logger, err = zapConfig.Build(zap.AddCaller(), zap.AddCallerSkip(1)) // Добавляем информацию о месте вызова
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}
}

func TestLoggerLevels(t *testing.T) {
	// Пример для уровня "DEBUG"
	t.Run("test debug level", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "DEBUG")
		InitLogger("DEBUG") // Инициализируем вручную

		// Замечаем вывод лога с помощью observer
		core, recordedLogs := observer.New(zap.DebugLevel)
		logger = zap.New(core)

		// Логируем сообщение
		Debug("Debug message")

		// Проверяем, что запись в логе действительно была
		logs := recordedLogs.TakeAll()
		require.Len(t, logs, 1, "Expected one log message, but got %d", len(logs))

		// Проверяем сообщение и уровень
		assert.Equal(t, "Debug message", logs[0].Message)
		assert.Equal(t, zap.DebugLevel, logs[0].Level)
	})

	// Пример для уровня "ERROR"
	t.Run("test error level", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "ERROR")
		InitLogger("ERROR")

		// Замечаем вывод лога с помощью observer
		core, recordedLogs := observer.New(zap.ErrorLevel)
		logger = zap.New(core)

		// Логируем сообщение
		Error("Error message")

		// Проверяем, что запись в логе действительно была
		logs := recordedLogs.TakeAll()
		require.Len(t, logs, 1, "Expected one log message, but got %d", len(logs))

		// Проверяем сообщение и уровень
		assert.Equal(t, "Error message", logs[0].Message)
		assert.Equal(t, zap.ErrorLevel, logs[0].Level)
	})

	// Пример для уровня по умолчанию "INFO"
	t.Run("test default info level", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "")
		InitLogger("") // Или можно передать "INFO" напрямую

		// Замечаем вывод лога с помощью observer
		core, recordedLogs := observer.New(zap.InfoLevel)
		logger = zap.New(core)

		// Логируем сообщение
		Info("Info message")

		// Проверяем, что запись в логе действительно была
		logs := recordedLogs.TakeAll()
		require.Len(t, logs, 1, "Expected one log message, but got %d", len(logs))

		// Проверяем сообщение и уровень
		assert.Equal(t, "Info message", logs[0].Message)
		assert.Equal(t, zap.InfoLevel, logs[0].Level)
	})

	// Пример для ERROR уровня с сообщением DEBUG (не должно быть записано)
	t.Run("test error level ignores debug message", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "ERROR")
		InitLogger("ERROR")

		// Замечаем вывод лога с помощью observer
		core, recordedLogs := observer.New(zap.ErrorLevel)
		logger = zap.New(core)

		// Логируем сообщение DEBUG, которое не должно быть записано
		Debug("Debug message")

		// Проверяем, что логов нет
		logs := recordedLogs.TakeAll()
		require.Len(t, logs, 0, "Expected no logs, but got %d", len(logs))
	})

	// Проверяем что для уровня DEBUG все сообщения должны быть записаны
	t.Run("test debug level logs all levels", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "DEBUG")
		InitLogger("DEBUG")

		// Замечаем вывод лога с помощью observer
		core, recordedLogs := observer.New(zap.DebugLevel)
		logger = zap.New(core)

		// Логируем сообщения разных уровней
		Debug("Debug message")
		Info("Info message")
		Error("Error message")

		// Проверяем, что все сообщения записались
		logs := recordedLogs.TakeAll()
		require.Len(t, logs, 3, "Expected three log messages, but got %d", len(logs))

		// Проверяем, что уровни и сообщения соответствуют
		assert.Equal(t, "Debug message", logs[0].Message)
		assert.Equal(t, zap.DebugLevel, logs[0].Level)

		assert.Equal(t, "Info message", logs[1].Message)
		assert.Equal(t, zap.InfoLevel, logs[1].Level)

		assert.Equal(t, "Error message", logs[2].Message)
		assert.Equal(t, zap.ErrorLevel, logs[2].Level)
	})

	// Проверим, что для уровня ERROR только ошибки записываются
	t.Run("test error level ignores info message", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "ERROR")
		InitLogger("ERROR")

		// Замечаем вывод лога с помощью observer
		core, recordedLogs := observer.New(zap.ErrorLevel)
		logger = zap.New(core)

		// Логируем сообщение INFO, которое не должно быть записано на уровне ERROR
		Info("This should not appear in error level logs")

		// Проверяем, что логов нет
		logs := recordedLogs.TakeAll()
		require.Len(t, logs, 0, "Expected no logs, but got %d", len(logs))
	})
}
