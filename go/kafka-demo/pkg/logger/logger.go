package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger 简单的日志记录器
type Logger struct {
	prefix string
	logger *log.Logger
}

// New 创建新的日志记录器
func New(prefix string) *Logger {
	return &Logger{
		prefix: prefix,
		logger: log.New(os.Stdout, "", 0),
	}
}

// Info 记录信息日志
func (l *Logger) Info(format string, v ...interface{}) {
	l.log("INFO", format, v...)
}

// Error 记录错误日志
func (l *Logger) Error(format string, v ...interface{}) {
	l.log("ERROR", format, v...)
}

// Warn 记录警告日志
func (l *Logger) Warn(format string, v ...interface{}) {
	l.log("WARN", format, v...)
}

// Debug 记录调试日志
func (l *Logger) Debug(format string, v ...interface{}) {
	l.log("DEBUG", format, v...)
}

func (l *Logger) log(level, format string, v ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, v...)
	l.logger.Printf("[%s] [%s] [%s] %s", timestamp, level, l.prefix, message)
}
