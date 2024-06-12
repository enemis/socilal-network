package logger

import (
	"github.com/sirupsen/logrus"
)

type Fields = map[string]any

type LoggerInterface interface {
	Debug(msg string, fields Fields)
	Info(msg string, fields Fields)
	Warn(msg string, fields Fields)
	Error(msg string, err error, fields Fields)
	Fatal(msg string, err error, fields Fields)
}

type Logger struct {
	logger *logrus.Logger
}

func NewLogger() *Logger {
	logger := logrus.New()
	logger.Formatter = new(logrus.JSONFormatter)

	return &Logger{
		logger: logger,
	}
}

func (s Logger) Debug(msg string, fields Fields) {
	s.logger.WithFields(fields).Debug(msg)
}

func (s Logger) Info(msg string, fields Fields) {
	s.logger.WithFields(fields).Info(msg)
}

func (s Logger) Warn(msg string, fields Fields) {
	s.logger.WithFields(fields).Warn(msg)
}

func (s Logger) Error(msg string, err error, fields Fields) {
	if fields == nil {
		fields = make(Fields, 1)
	}

	fields["error"] = err

	s.logger.WithFields(fields).Error(msg)
}

func (s Logger) Fatal(msg string, err error, fields Fields) {
	if fields == nil {
		fields = make(map[string]interface{}, 1)
	}

	fields["error"] = err

	s.logger.Fatal(msg)
}
