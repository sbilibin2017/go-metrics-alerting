package logger

// LogLevel is a custom type for defining log levels.
type LogLevel int

// Log levels
const (
	DEBUG LogLevel = iota
	INFO
	ERROR
)

// LoggerConfig holds the configuration settings for the logger.
type LoggerConfig struct {
	LogLevel LogLevel
}
