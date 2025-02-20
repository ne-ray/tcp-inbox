package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Interface -.
type Interface interface {
	With(field string, value interface{}) InterfaceLogger
	Debug(message string)
	Info(message string)
	Warn(message string)
	Warning(message string)
	Error(message string)
	Fatal(message string)
	Shutdown() error
}

// InterfaceLogger -.
type InterfaceLogger interface {
	With(field string, value interface{}) InterfaceLogger
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
}

type LoggerInstance struct {
	logger *zap.SugaredLogger
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
			MessageKey: "message",
			LevelKey:   "level",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		},
	}

	logger := zap.Must(cfg.Build())

	return &Logger{
		logger:     logger,
		appName:    appName,
		appVersion: appVer,
	}
}

func (l *Logger) newInstance() InterfaceLogger {
	return &LoggerInstance{
		logger: l.logger.Sugar(),
	}
}

func (l *Logger) With(field string, value interface{}) InterfaceLogger {
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
func (l *LoggerInstance) With(field string, value interface{}) InterfaceLogger {
	l.logger.With(field, value)

	return l
}

// Debug -.
func (l *LoggerInstance) Debug(message string) {
	l.logger.Debug(message)
}

// Info -.
func (l *LoggerInstance) Info(message string) {
	l.logger.Info(message)
}

// Warn -.
func (l *LoggerInstance) Warn(message string) {
	l.logger.Warn(message)
}

// Warning -.
func (l *LoggerInstance) Warning(message string) {
	l.logger.Warn(message)
}

// Error -.
func (l *LoggerInstance) Error(message string) {
	l.logger.Error(message)
}

// Fatal -.
func (l *LoggerInstance) Fatal(message string) {
	l.logger.Fatal(message)

	os.Exit(1)
}
