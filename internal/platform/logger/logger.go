package logger

import (
	"encoding/json"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/platform/config"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

type Field struct {
	Key   string
	Value any
}

type Logger struct {
	output io.Writer
	level  Level
	mutex  sync.Mutex
	now    func() time.Time
}

func New(output io.Writer, level config.LogLevel) *Logger {
	return &Logger{
		output: output,
		level:  parseLevel(level),
		now:    time.Now,
	}
}

func (log *Logger) Debug(message string, fields ...Field) {
	log.write(DebugLevel, "debug", message, fields...)
}

func (log *Logger) Info(message string, fields ...Field) {
	log.write(InfoLevel, "info", message, fields...)
}

func (log *Logger) Warn(message string, fields ...Field) {
	log.write(WarnLevel, "warn", message, fields...)
}

func (log *Logger) Error(message string, fields ...Field) {
	log.write(ErrorLevel, "error", message, fields...)
}

func String(key string, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func (log *Logger) write(level Level, levelName string, message string, fields ...Field) {
	if level < log.level {
		return
	}

	entry := map[string]any{
		"level":      levelName,
		"message":    message,
		"created_at": log.now().UTC().Format(time.RFC3339Nano),
	}
	for _, field := range fields {
		if field.Key == "" {
			continue
		}
		entry[field.Key] = redactValue(field.Key, field.Value)
	}

	log.mutex.Lock()
	defer log.mutex.Unlock()
	_ = json.NewEncoder(log.output).Encode(entry)
}

func parseLevel(level config.LogLevel) Level {
	switch level {
	case config.LogLevelDebug:
		return DebugLevel
	case config.LogLevelWarn:
		return WarnLevel
	case config.LogLevelError:
		return ErrorLevel
	default:
		return InfoLevel
	}
}

func redactValue(key string, value any) any {
	normalized := strings.ToLower(key)
	secretFragments := []string{
		"password",
		"token",
		"api_key",
		"apikey",
		"secret",
		"private_key",
		"credential",
		"authorization",
		"cookie",
	}

	for _, fragment := range secretFragments {
		if strings.Contains(normalized, fragment) {
			return "[REDACTED]"
		}
	}
	return value
}
