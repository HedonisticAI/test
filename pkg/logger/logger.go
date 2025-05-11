package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	Log *logrus.Logger
}

func NewLogger() *Logger {
	Log := logrus.New()
	Log.SetLevel(logrus.DebugLevel)
	return &Logger{Log: Log}
}

func (L *Logger) Debug(Msg ...interface{}) {
	L.Log.Debug(Msg...)
}

func (L *Logger) Error(Msg ...interface{}) {
	L.Log.Error(Msg...)
}

func (L *Logger) Info(Msg ...interface{}) {
	L.Log.Info(Msg...)
}
