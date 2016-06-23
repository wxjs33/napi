package log

import (
	"code.google.com/p/log4go"
	"github.com/wxjs33/napi/variable"
)

type Log struct {
	l log4go.Logger
}

func GetLogger(path string, level string) *Log {
	var log Log

	if path == "" {
		path = variable.DEFAULT_LOG_PATH
	}

	lv := log4go.ERROR
	switch level {
	case "debug":
		lv = log4go.DEBUG
	case "error":
		lv = log4go.ERROR
	case "info":
		lv = log4go.INFO
	}

	l := log4go.NewDefaultLogger(lv)
	flw := log4go.NewFileLogWriter(path, false)
	if flw == nil {
		return nil
	}
	flw.SetFormat("[%D %T] [%L] %M")
	//flw.SetRotate(true)
	//flw.SetRotateLines(50)
	flw.SetRotateDaily(true)
	l.AddFilter("log", lv, flw)

	log.l = l

	return &log
}

func (l *Log) Info(arg0 interface{}, args ...interface{}) {
	l.l.Info(arg0, args...)
}

func (l *Log) Error(arg0 interface{}, args ...interface{}) {
	l.l.Error(arg0, args...)
}

func (l *Log) Debug(arg0 interface{}, args ...interface{}) {
	l.l.Debug(arg0, args...)
}
