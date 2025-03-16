package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Severity string

const (
	SeverityInfo  Severity = "INFO"
	SeverityWarn  Severity = "WARN"
	SeverityDebug Severity = "DEBUG"
	SeverityError Severity = "ERROR"
	SeverityFatal Severity = "FATAL"
)

func defaultConfig() *Config {
	return &Config{
		Output: "stdout",
		Name:   "LOGGER",
		Level:  SeverityInfo,
	}
}

type Logger struct {
	cnf *Config
}

type Config struct {
	Name   string
	Level  Severity
	Output string // stdout, <filepath>
}

func NewLogger(config *Config) *Logger {
	var logger = new(Logger)
	if config != nil {
		if config.Output == "" {
			config.Output = defaultConfig().Output
		}
		if config.Name == "" {
			config.Name = defaultConfig().Name
		}
		if string(config.Level) == "" {
			config.Level = defaultConfig().Level
		}
	} else {
		config = defaultConfig()
	}

	logger.cnf = config

	if logger.cnf.Output != "stdout" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   logger.cnf.Output,
			MaxSize:    5, // megabytes
			MaxBackups: 100,
			MaxAge:     30, // days
		})

		log.SetFlags(0)
	}

	return logger
}

func (l Logger) LogF(level Severity, format string, args ...any) string {
	message := fmt.Sprintf(
		"%s %s %s: %s",
		time.Now().UTC().Format(time.RFC3339Nano),
		l.cnf.Name, level, fmt.Sprint(format, args),
	)

	switch l.cnf.Output {
	case "stdout":
		fmt.Println(message)
	default:
		log.Println(message)
	}

	return message
}

func (l Logger) ErrorF(format string, args ...any) string {
	return l.LogF(SeverityError, format, args...)
}

func (l Logger) Fatal(format string, args ...any) {
	l.LogF(SeverityFatal, format, args...)
	os.Exit(1)
}

func (l Logger) Info(format string, args ...any) string {
	return l.LogF(SeverityInfo, format, args...)
}

func (l Logger) Debug(format string, args ...any) string {
	return l.LogF(SeverityDebug, format, args...)
}

func (l Logger) Warn(format string, args ...any) string {
	return l.LogF(SeverityWarn, format, args...)
}
