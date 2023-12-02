package watermillkratos

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/go-kratos/kratos/v2/log"
)

var _ watermill.LoggerAdapter = (*Logger)(nil)

const (
	errorkey = "error"
)

type Logger struct {
	msgKey string
	logger *log.Helper
	fields watermill.LogFields
}

func NewLogger(
	logger *log.Helper,
	msgKey string,
) *Logger {
	return &Logger{
		msgKey: msgKey,
		logger: logger,
		fields: make(watermill.LogFields),
	}
}

func toKratosLogFields(fields watermill.LogFields) []interface{} {
	kratosFields := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		kratosFields = append(kratosFields, k, v)
	}
	return kratosFields
}

func (l *Logger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &Logger{
		logger: l.logger,
		fields: l.fields.Add(fields),
	}
}

func (l *Logger) Error(msg string, err error, fields watermill.LogFields) {
	fields = l.fields.Add(fields).Add(watermill.LogFields{
		l.msgKey: msg,
	})
	fields = fields.Add(watermill.LogFields{
		errorkey: err,
	})
	l.logger.Errorw(toKratosLogFields(fields)...)
}

func (l *Logger) Info(msg string, fields watermill.LogFields) {
	fields = l.fields.Add(fields).Add(watermill.LogFields{
		l.msgKey: msg,
	})
	l.logger.Infow(toKratosLogFields(fields)...)
}

func (l *Logger) Debug(msg string, fields watermill.LogFields) {
	fields = l.fields.Add(fields).Add(watermill.LogFields{
		l.msgKey: msg,
	})
	l.logger.Debugw(toKratosLogFields(fields)...)
}

func (l *Logger) Trace(msg string, fields watermill.LogFields) {
	fields = l.fields.Add(fields).Add(watermill.LogFields{
		l.msgKey: msg,
	})
	l.logger.Debugw(toKratosLogFields(fields)...)
}
