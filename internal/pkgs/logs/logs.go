package logs

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log    *zap.Logger
	config zap.Config
	err    error
)

func LogInit() {
	if viper.GetString("server.mode") == "debug" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = ""

	log, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

func Info(message string, field ...zapcore.Field) {
	log.Info(message, field...)
}

func Debug(message string, field ...zapcore.Field) {
	log.Debug(message, field...)
}
func Error(message interface{}, field ...zapcore.Field) {
	switch v := message.(type) {
	case error:
		log.Error(v.Error(), field...)
	case string:
		log.Error(v, field...)
	}
}
