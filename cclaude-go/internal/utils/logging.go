package utils

import (
	"fmt"
	"time"
)

type LogLevel int

const (
	LevelSilent LogLevel = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

var currentLevel = LevelInfo

func SetLogLevel(level LogLevel) {
	currentLevel = level
}

func log(level LogLevel, format string, args ...interface{}) {
	if level > currentLevel {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	prefix := ""

	switch level {
	case LevelError:
		prefix = "❌"
	case LevelWarn:
		prefix = "⚠️"
	case LevelInfo:
		prefix = "ℹ️"
	case LevelDebug:
		prefix = "🔍"
	}

	message := fmt.Sprintf(format, args...)
	fmt.Printf("[%s] %s %s\n", timestamp, prefix, message)
}

func Info(format string, args ...interface{})  { log(LevelInfo, format, args...) }
func Warn(format string, args ...interface{})  { log(LevelWarn, format, args...) }
func Error(format string, args ...interface{}) { log(LevelError, format, args...) }
func Debug(format string, args ...interface{}) { log(LevelDebug, format, args...) }