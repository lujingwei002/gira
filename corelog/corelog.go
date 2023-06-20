package corelog

import (
	"os"
	"path/filepath"

	"github.com/lujingwei002/gira"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger Logger

type Logger interface {
	Infow(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Fatal(args ...interface{})
	Warn(args ...interface{})
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Debugf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	Warnf(template string, args ...interface{})
}

func ConfigLog(config gira.LogConfig) error {
	cores := make([]zapcore.Core, 0)
	// 1.控制台输出
	if config.Console {
		var level zap.AtomicLevel
		err := level.UnmarshalText([]byte(config.Level))
		if err != nil {
			return err
		}
		enabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= level.Level()
		})
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig), // 控制台输出格式
			zapcore.AddSync(os.Stdout),               // 输出到控制台
			enabler,
		)
		cores = append(cores, consoleCore)
	}

	// 2.滚动文件输出
	if config.File {
		var level zap.AtomicLevel
		err := level.UnmarshalText([]byte(config.Level))
		if err != nil {
			return err
		}
		enabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= level.Level()
		})

		encoderCfg := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		rollingCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg), // 滚动日志输出格式
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   filepath.Join(config.Dir, config.Name), // 日志文件路径
				MaxSize:    config.MaxSize,                         // 每个日志文件的最大大小，单位为 MB
				MaxBackups: config.MaxBackups,                      // 保留的旧日志文件的最大个数
				MaxAge:     config.MaxAge,                          // 保留的旧日志文件的最大天数
				Compress:   config.Compress,                        // 是否压缩旧日志文件
			}),
			enabler,
		)
		cores = append(cores, rollingCore)
	}
	// 创建日志对象
	logger := zap.New(zapcore.NewTee(cores...))
	logger = logger.WithOptions(zap.WithCaller(true), zap.AddCallerSkip(1))
	sugar := logger.Sugar()
	defaultLogger = sugar
	return nil
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
	defaultLogger.Debug(args...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

func Println(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
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
func Warnw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Warnw(msg, keysAndValues...)
}

func Printf(template string, args ...interface{}) {
	defaultLogger.Infof(template, args...)
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

func Warnf(template string, args ...interface{}) {
	defaultLogger.Warnf(template, args...)
}
