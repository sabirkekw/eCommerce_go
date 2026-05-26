package logger

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func InitLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level.SetLevel(zap.DebugLevel)
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	Log = logger.Sugar()
}
