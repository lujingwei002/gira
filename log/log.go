package log

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lujingwei002/gira"
	"github.com/lujingwei002/gira/facade"
	"github.com/lujingwei002/gira/proj"
	"github.com/natefinch/lumberjack"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func ConfigCliLog() error {
	// 配置日志输出
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
		zap.NewAtomicLevelAt(zap.DebugLevel),
	)
	// 滚动日志配置
	rollingCfg := zapcore.EncoderConfig{
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
		zapcore.NewConsoleEncoder(rollingCfg), // 滚动日志输出格式
		// zapcore.NewJSONEncoder(rollingCfg), // 滚动日志输出格式
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(proj.Config.LogDir, fmt.Sprintf("cli.log")), // 日志文件路径
			MaxSize:    10 * 1024,                                                 // 每个日志文件的最大大小，单位为 MB
			MaxBackups: 10,                                                        // 保留的旧日志文件的最大个数
			MaxAge:     30,                                                        // 保留的旧日志文件的最大天数
			Compress:   true,                                                      // 是否压缩旧日志文件
		}),
		zap.NewAtomicLevelAt(zap.DebugLevel),
	)

	// 创建日志对象
	logger := zap.New(zapcore.NewTee(consoleCore, rollingCore))
	logger = logger.WithOptions(zap.WithCaller(true), zap.AddCallerSkip(1))
	sugar := logger.Sugar()
	defaultLogger = sugar
	return nil
}

type MongoSink struct {
	collection  string
	writeOption *options.InsertOneOptions
}

func (s *MongoSink) Write(data []byte) (n int, err error) {
	client := facade.GetLogDbClient()
	if client == nil {
		return
	}
	switch c := client.(type) {
	case gira.MongoClient:
		doc := bson.M{}
		if err = json.Unmarshal(data, &doc); err != nil {
			return
		}
		doc["app_id"] = facade.GetAppId()
		doc["app_full_name"] = facade.GetAppFullName()
		database := c.GetMongoDatabase()
		_, err = database.Collection(s.collection).InsertOne(context.Background(), doc, s.writeOption)
		return
	default:
		err = gira.ErrDbNotSupport.Trace()
		return
	}
}

func ConfigLog(facade gira.Application, config gira.LogConfig) error {
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
				Filename:   filepath.Join(facade.GetLogDir(), fmt.Sprintf("%s.log", facade.GetAppFullName())), // 日志文件路径
				MaxSize:    config.MaxSize,                                                                    // 每个日志文件的最大大小，单位为 MB
				MaxBackups: config.MaxBackups,                                                                 // 保留的旧日志文件的最大个数
				MaxAge:     config.MaxAge,                                                                     // 保留的旧日志文件的最大天数
				Compress:   config.Compress,                                                                   // 是否压缩旧日志文件
			}),
			enabler,
		)
		cores = append(cores, rollingCore)
	}
	// 3.数据库输出
	if config.Db {
		var level zap.AtomicLevel
		err := level.UnmarshalText([]byte(config.DbLevel))
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
		mongoCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg), // 滚动日志输出格式
			zapcore.AddSync(&MongoSink{
				collection: "log",
				writeOption: options.InsertOne().
					SetBypassDocumentValidation(true),
			}),
			enabler,
		)
		cores = append(cores, mongoCore)
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
	defaultLogger.Info(args...)
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
