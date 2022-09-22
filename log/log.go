package log

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func Init(config *Config) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	cores := make([]zapcore.Core, 1)
	cores[0] = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)

	if config.Filename != "" {
		priority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev >= zapcore.Level(config.Level)
		})
		infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    2,
			MaxBackups: 100,
			MaxAge:     30,
			Compress:   false,
		})
		infoFileCore := zapcore.NewCore(encoder, infoFileWriteSyncer, priority)
		cores = append(cores, infoFileCore)
	}

	logger = zap.New(zapcore.NewTee(cores...), zap.AddCaller()).Sugar()
}

func Info(args ...interface{}) {
	logger.Info(args)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args)
}

func Error(args ...interface{}) {
	logger.Error(args)
}
