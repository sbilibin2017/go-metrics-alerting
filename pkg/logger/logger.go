// pkg/logger/logger.go
package logger

import (
	"github.com/sirupsen/logrus"
)

// Logger представляет собой глобальный инстанс для удобства использования
var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetLevel(logrus.DebugLevel)
}
