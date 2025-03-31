package logger

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const timeFormat = time.RFC3339

// Interface -.
type Interface interface {
	InterfaceLogger
	Shutdown() error
}

type InterfaceLogger interface {
	With(field string, value any) InterfaceLogger
	Debug(message string)
	Info(message string)
	Warn(message string)
	Warning(message string)
	Error(message string)
	Fatal(message string)
}

// Logger -.
type Logger struct {
	logger     *zap.Logger
	appName    string
	appVersion string
	timeFormat string
}

type LoggerInstance struct {
	logger     *zap.SugaredLogger
	timeFormat string
}

var _ Interface = (*Logger)(nil)

// New -.
func New(level, appName, appVer string) *Logger {
	var l zapcore.Level

	switch strings.ToLower(level) {
	case "error":
		l = zap.ErrorLevel
	case "warn":
		l = zap.WarnLevel
	case "info":
		l = zap.InfoLevel
	case "debug":
		l = zap.DebugLevel
	default:
		l = zap.InfoLevel
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(l),
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    map[string]interface{}{"application_name": appName, "application_version": appVer},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		},
	}

	logger := zap.Must(cfg.Build())

	return &Logger{
		logger:     logger,
		appName:    appName,
		appVersion: appVer,
		timeFormat: timeFormat,
	}
}

func (l *Logger) newInstance() InterfaceLogger {
	return &LoggerInstance{
		logger:     l.logger.Sugar(),
		timeFormat: l.timeFormat,
	}
}

func (l *Logger) With(field string, value any) InterfaceLogger {
	return l.newInstance().With(field, value)
}

// Debug -.
func (l *Logger) Debug(message string) {
	l.newInstance().Debug(message)
}

// Info -.
func (l *Logger) Info(message string) {
	l.newInstance().Info(message)
}

// Warn -.
func (l *Logger) Warn(message string) {
	l.newInstance().Warn(message)
}

// Warning -.
func (l *Logger) Warning(message string) {
	l.newInstance().Warning(message)
}

// Error -.
func (l *Logger) Error(message string) {
	l.newInstance().Error(message)
}

// Fatal -.
func (l *Logger) Fatal(message string) {
	l.newInstance().Fatal(message)
}

func (l *Logger) Shutdown() error {
	return l.logger.Sync()
}

// With -.
func (li *LoggerInstance) With(field string, value any) InterfaceLogger {
	li.logger = li.logger.With(field, value)

	return li
}

// WithTime -.
func (li *LoggerInstance) WithTime() InterfaceLogger {
	li.With("time", time.Now().Format(li.timeFormat))

	return li
}

// Debug -.
func (li *LoggerInstance) Debug(message string) {
	li.WithTime()

	li.logger.Debug(message)
}

// Info -.
func (li *LoggerInstance) Info(message string) {
	li.WithTime()

	li.logger.Info(message)
}

// Warn -.
func (li *LoggerInstance) Warn(message string) {
	li.WithTime()

	li.logger.Warn(message)
}

// Warning -.
func (li *LoggerInstance) Warning(message string) {
	li.WithTime()

	li.logger.Warn(message)
}

// Error -.
func (li *LoggerInstance) Error(message string) {
	li.WithTime()

	li.logger.Error(message)
}

// Fatal -.
func (li *LoggerInstance) Fatal(message string) {
	li.WithTime()

	li.logger.Fatal(message)

	os.Exit(1)
}
