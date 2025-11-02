package logger

import (
	"market/pkg/config"

	"go.uber.org/zap"
)

func NewLogger() *zap.SugaredLogger {
	var logger *zap.Logger
	var err error

	if config.Get().ENV == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic("Failed to create logger: " + err.Error())
	}

	if logger == nil {
		panic("Logger is nil")
	}

	defer logger.Sync()

	return logger.Sugar()
}
