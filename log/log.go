package log

import (
	"go.uber.org/zap"
)

var defaultLogger Logger

type Logger interface {
	Infow(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Fatal(args ...interface{})
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Debugf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

func init() {
	logger, _ := zap.NewDevelopment()
	logger = logger.WithOptions(zap.AddCallerSkip(1))
	sugar := logger.Sugar()
	defaultLogger = sugar
}

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Debug(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Error(args ...interface{}) {
	defaultLogger.Info(args...)
}
func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}
func Infow(msg string, keysAndValues ...interface{}) {
	defaultLogger.Infow(msg, keysAndValues...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Debugw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Errorw(msg, keysAndValues...)
}
func Fatalw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Fatalw(msg, keysAndValues...)
}
func Infof(template string, args ...interface{}) {
	defaultLogger.Infof(template, args...)
}

func Debugf(template string, args ...interface{}) {
	defaultLogger.Debugf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	defaultLogger.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	defaultLogger.Fatalf(template, args...)
}
