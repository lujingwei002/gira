package gira

import (
	"go.uber.org/zap"
)

var defaultLogger Logger

type Logger interface {
	Infow(msg string, keysAndValues ...interface{})
}

func init() {
	logger, _ := zap.NewDevelopment()
	logger = logger.WithOptions(zap.AddCallerSkip(1))
	sugar := logger.Sugar()
	defaultLogger = sugar

}

func Infow(msg string, keysAndValues ...interface{}) {
	defaultLogger.Infow(msg, keysAndValues...)
}
