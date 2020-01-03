package logger

import "go.uber.org/zap"

func SetupLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(logger)
	zap.RedirectStdLog(logger)
	return logger, nil
}
