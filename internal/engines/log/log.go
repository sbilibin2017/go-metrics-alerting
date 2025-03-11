package log

import (
	"fmt"
	"os"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	sugar  *zap.SugaredLogger // Sugared logger
	logger *zap.Logger        // Regular zap logger
)

// Debugw logs a message with the Debug level with structured fields using the Sugared Logger
func Debugw(msg string, keysAndValues ...interface{}) {
	sugar.Debugw(msg, keysAndValues...)
}

// Infow logs a message with the Info level with structured fields using the Sugared Logger
func Infow(msg string, keysAndValues ...interface{}) {
	sugar.Infow(msg, keysAndValues...)
}

// Errorw logs a message with the Error level with structured fields using the Sugared Logger
func Errorw(msg string, keysAndValues ...interface{}) {
	sugar.Errorw(msg, keysAndValues...)
}

// Debugw logs a message with the Debug level with structured fields using the Regular zap Logger
func Debug(msg string, args ...interface{}) {
	logger.Debug(msg, zap.String("message", fmt.Sprintf(msg, args...)))
}

// Info logs a message with the Info level using the Regular zap Logger
func Info(msg string, args ...interface{}) {
	logger.Info(msg, zap.String("message", fmt.Sprintf(msg, args...)))
}

// Error logs a message with the Error level using the Regular zap Logger
func Error(msg string, args ...interface{}) {
	logger.Error(msg, zap.String("message", fmt.Sprintf(msg, args...)))
}

// getCallerInfo retrieves the caller information (file and line number)
func getCallerInfo() string {
	_, file, line, _ := runtime.Caller(2)   // Get the caller info from two levels up
	return fmt.Sprintf("%s:%d", file, line) // Format it as "file:line"
}

func init() {
	// Create encoder configuration for both loggers
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "msg",    // Keep "msg" for the message field
		LevelKey:     "level",  // Use "level" for the log level
		TimeKey:      "ts",     // Use "ts" for timestamp
		NameKey:      "logger", // Keep "logger" for logger name (not used in this case)
		CallerKey:    "caller", // Keep "caller" for the file and line number
		LineEnding:   zapcore.DefaultLineEnding,
		EncodeTime:   zapcore.ISO8601TimeEncoder,  // Format the time in ISO8601 format
		EncodeLevel:  zapcore.CapitalLevelEncoder, // Capitalize the level (INFO, DEBUG)
		EncodeCaller: zapcore.ShortCallerEncoder,  // Shorten the caller field to just the filename:line number
	}

	// Create encoder and core for both loggers
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel)

	// Create the Sugared logger
	loggerSugared := zap.New(core)
	sugar = loggerSugared.Sugar()

	// Create the regular zap logger
	regularLogger := zap.New(core)
	logger = regularLogger
}
