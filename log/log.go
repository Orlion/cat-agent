package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.SugaredLogger

func Init(config *Config) {
	config = withDefaultConf(config)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	level, err := zapcore.ParseLevel(config.StdoutLevel)
	if err != nil {
		level = zapcore.InfoLevel
	}

	cores := make([]zapcore.Core, 1)
	cores[0] = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)

	level, err = zapcore.ParseLevel(config.StdoutLevel)
	if err != nil {
		level = zapcore.ErrorLevel
	}
	if config.Filename != "" {
		priority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev >= level
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

	logger = zap.New(zapcore.NewTee(cores...)).Sugar()
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func Shutdown() {
	logger.Info("log shutdown...")
	logger.Sync()
}
