package logger

import (
	"go.uber.org/zap"
)

type CrowLogger struct {
	*zap.SugaredLogger
}

var Global CrowLogger

func Init() {
	zapLogger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(zapLogger)
	Global = CrowLogger{
		zap.S(),
	}
}
