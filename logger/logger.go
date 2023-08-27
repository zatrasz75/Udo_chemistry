package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var logger *log.Logger

func init() {
	logFile, err := os.Create("app.log")
	if err != nil {
		log.Fatal("Ошибка создания файла журнала:", err)
	}
	logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds)
}

func logWithCallerInfo(file string, line int, level string, message string, args ...interface{}) {
	caller := fmt.Sprintf("%s:%d", filepath.Base(file), line)
	messageWithCaller := fmt.Sprintf("[%s] %s %s %s", level, getFormattedTime(), caller, fmt.Sprintf(message, args...))
	logger.Println(messageWithCaller)
	fmt.Println(messageWithCaller)
}

// Error записывает сообщение об ошибке в лог вместе с контекстом вызова функции.
// Параметр err содержит ошибку, связанную с данным сообщением.
func Error(message string, err error) {
	_, file, line, _ := runtime.Caller(1)
	logWithCallerInfo(file, line, "ERROR", "%s: %v", message, err)
}

// Info записывает информационное сообщение в лог вместе с контекстом вызова функции.
// Параметры args содержат дополнительные данные для сообщения.
func Info(message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logWithCallerInfo(file, line, "INFO", message, args...)
}

// Fatal записывает фатальное сообщение в лог вместе с контекстом вызова функции
// и завершает приложение с кодом ошибки 1.
// Параметр err содержит ошибку, связанную с данным сообщением.
func Fatal(message string, err error) {
	_, file, line, _ := runtime.Caller(1)
	logWithCallerInfo(file, line, "FATAL", "%s: %v", message, err)
	os.Exit(1) // Завершаем приложение с кодом ошибки
}

func getFormattedTime() string {
	return time.Now().Format("2006-01-02 15:04:05.000000")
}
