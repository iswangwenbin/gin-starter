package logx

import (
	"sync"

	"go.uber.org/zap"
)

var (
	sugar   *zap.SugaredLogger
	once    sync.Once
	initErr error
)

func ensure() error {
	once.Do(func() {
		var logger *zap.Logger
		// 这里可以选择 NewProduction() 或 NewDevelopment()
		logger, initErr = zap.NewProduction()
		if initErr != nil {
			return
		}
		zap.ReplaceGlobals(logger)
		sugar = logger.Sugar()
	})
	return initErr
}

func GetLogger() (*zap.SugaredLogger, error) {
	if err := ensure(); err != nil {
		return nil, err
	}
	return sugar, nil
}

func Info(msg string, keysAndValues ...interface{}) {
	if err := ensure(); err != nil {
		// 无法初始化时，可选择 panic 或者忽略
		panic("logx initialization failed: " + err.Error())
	}
	sugar.Infow(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	if err := ensure(); err != nil {
		panic("logx initialization failed: " + err.Error())
	}
	sugar.Errorw(msg, keysAndValues...)
}

func Sync() error {
	if err := ensure(); err != nil {
		return err
	}
	return sugar.Sync()
}
