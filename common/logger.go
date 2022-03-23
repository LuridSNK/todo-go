package common

import (
	"go.uber.org/zap"
)

type AppLogger struct {
	*zap.SugaredLogger
}

func NewLogger(env string) (*AppLogger, error) {
	if env == "development" {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
		return &AppLogger{logger.Sugar()}, nil
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &AppLogger{logger.Sugar()}, nil
}
