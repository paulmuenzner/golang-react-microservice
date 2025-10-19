package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

var (
	debugLogger = log.New(os.Stdout, "", 0)
	infoLogger  = log.New(os.Stdout, "", 0)
	warnLogger  = log.New(os.Stdout, "", 0)
	errorLogger = log.New(os.Stderr, "", 0)
	fatalLogger = log.New(os.Stderr, "", 0)
)

// SetLoggerWriters allows tests to redirect logger outputs.
func SetLoggerWriters(outWriter, errWriter *os.File) {
	debugLogger.SetOutput(outWriter)
	infoLogger.SetOutput(outWriter)
	warnLogger.SetOutput(outWriter)
	errorLogger.SetOutput(errWriter)
	fatalLogger.SetOutput(errWriter)
}

// formatMessage builds a detailed log message with timestamp, level, file, and line number
func formatMessage(level string, msg string) string {
	// Get caller info (skip 2 levels: this function and the public log function)
	_, file, line, ok := runtime.Caller(2)
	location := "unknown"
	if ok {
		location = fmt.Sprintf("%s:%d", file, line)
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	return fmt.Sprintf("[%s] [%s] [%s] %s", timestamp, level, location, msg)
}

// Debug logs detailed developer-level information
func Debug(msg string) {
	debugLogger.Println(formatMessage("DEBUG", msg))
}

// Info logs general application flow messages
func Info(msg string) {
	infoLogger.Println(formatMessage("INFO", msg))
}

// Warn logs potential issues or unexpected states
func Warn(msg string) {
	warnLogger.Println(formatMessage("WARN", msg))
}

// Error logs actual errors that do not stop execution
func Error(msg string) {
	errorLogger.Println(formatMessage("ERROR", msg))
}

// Fatal logs a critical error and exits the application
func Fatal(msg string) {
	fatalLogger.Println(formatMessage("FATAL", msg))
	os.Exit(1)
}
