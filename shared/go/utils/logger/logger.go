package logger

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger is the application logger with context support
type Logger struct {
	zlog zerolog.Logger
}

var AppLogger *Logger

// Init initializes the global logger
// serviceName: Name of the microservice (e.g. "gateway", "service-a")
// environment: "development" or "production"
func Init(serviceName, environment string) {
	var output = os.Stdout

	// Development: Pretty console output with colors
	if environment == "development" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}
	}

	// Set log level based on environment
	level := zerolog.InfoLevel
	if environment == "development" {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)

	// Create logger with service name
	logger := zerolog.New(output).
		With().
		Timestamp().
		Str("service", serviceName).
		Logger()

	AppLogger = &Logger{zlog: logger}
	log.Logger = logger // Set as global default
}

// WithContext creates a logger with context values (Request ID, User ID, etc.)
func (l *Logger) WithContext(ctx context.Context) *Logger {
	newLogger := l.zlog.With()

	// Extract common context values
	if reqID := ctx.Value("request_id"); reqID != nil {
		newLogger = newLogger.Str("request_id", reqID.(string))
	}
	if userID := ctx.Value("user_id"); userID != nil {
		newLogger = newLogger.Str("user_id", userID.(string))
	}

	return &Logger{zlog: newLogger.Logger()}
}

// WithField adds a single field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		zlog: l.zlog.With().Interface(key, value).Logger(),
	}
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	ctx := l.zlog.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return &Logger{zlog: ctx.Logger()}
}

// Debug logs debug-level messages (only in development)
func (l *Logger) Debug(msg string) {
	l.zlog.Debug().Msg(msg)
}

// DebugWithFields logs debug with additional context
func (l *Logger) DebugWithFields(msg string, fields map[string]interface{}) {
	event := l.zlog.Debug()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Info logs general information
func (l *Logger) Info(msg string) {
	l.zlog.Info().Msg(msg)
}

// InfoWithFields logs info with additional context
func (l *Logger) InfoWithFields(msg string, fields map[string]interface{}) {
	event := l.zlog.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Warn logs warnings
func (l *Logger) Warn(msg string) {
	l.zlog.Warn().Msg(msg)
}

// WarnWithFields logs warning with additional context
func (l *Logger) WarnWithFields(msg string, fields map[string]interface{}) {
	event := l.zlog.Warn()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Error logs errors with error object
func (l *Logger) Error(msg string, err error) {
	l.zlog.Error().Err(err).Msg(msg)
}

// ErrorWithFields logs error with additional context
func (l *Logger) ErrorWithFields(msg string, err error, fields map[string]interface{}) {
	event := l.zlog.Error().Err(err)
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Fatal logs fatal error and exits application
func (l *Logger) Fatal(msg string, err error) {
	l.zlog.Fatal().Err(err).Msg(msg)
}

// HTTP logs HTTP requests with standard fields
func (l *Logger) HTTP(method, path string, statusCode int, duration time.Duration, clientIP string) {
	l.zlog.Info().
		Str("method", method).
		Str("path", path).
		Int("status", statusCode).
		Dur("duration_ms", duration).
		Str("client_ip", clientIP).
		Msg("HTTP Request")
}

// HTTPWithFields logs HTTP request with custom fields
func (l *Logger) HTTPWithFields(method, path string, statusCode int, duration time.Duration, fields map[string]interface{}) {
	event := l.zlog.Info().
		Str("method", method).
		Str("path", path).
		Int("status", statusCode).
		Dur("duration_ms", duration)

	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg("HTTP Request")
}

// ==========================================
// Package-level convenience functions
// (Can be used without creating Logger instance)
// ==========================================

func Debug(msg string) {
	if AppLogger != nil {
		AppLogger.Debug(msg)
	}
}

func DebugWithFields(msg string, fields map[string]interface{}) {
	if AppLogger != nil {
		AppLogger.DebugWithFields(msg, fields)
	}
}

func Info(msg string) {
	if AppLogger != nil {
		AppLogger.Info(msg)
	}
}

func InfoWithFields(msg string, fields map[string]interface{}) {
	if AppLogger != nil {
		AppLogger.InfoWithFields(msg, fields)
	}
}

func Warn(msg string) {
	if AppLogger != nil {
		AppLogger.Warn(msg)
	}
}

func WarnWithFields(msg string, fields map[string]interface{}) {
	if AppLogger != nil {
		AppLogger.WarnWithFields(msg, fields)
	}
}

func Error(msg string, err error) {
	if AppLogger != nil {
		AppLogger.Error(msg, err)
	}
}

func ErrorWithFields(msg string, err error, fields map[string]interface{}) {
	if AppLogger != nil {
		AppLogger.ErrorWithFields(msg, err, fields)
	}
}

func Fatal(msg string, err error) {
	if AppLogger != nil {
		AppLogger.Fatal(msg, err)
	}
}

func HTTP(method, path string, statusCode int, duration time.Duration, clientIP string) {
	if AppLogger != nil {
		AppLogger.HTTP(method, path, statusCode, duration, clientIP)
	}
}

// WithRequestID creates a logger with request ID
func WithRequestID(requestID string) *Logger {
	if AppLogger != nil {
		return AppLogger.WithField("request_id", requestID)
	}
	return nil
}

// WithUserID creates a logger with user ID
func WithUserID(userID string) *Logger {
	if AppLogger != nil {
		return AppLogger.WithField("user_id", userID)
	}
	return nil
}
