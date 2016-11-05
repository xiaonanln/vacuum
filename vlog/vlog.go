package vlog

import (
	"runtime/debug"

	sublog "github.com/Sirupsen/logrus"
)

var (
	DEBUG = sublog.DebugLevel
	INFO  = sublog.InfoLevel
	WARN  = sublog.WarnLevel
	ERROR = sublog.ErrorLevel
	PANIC = sublog.PanicLevel
	FATAL = sublog.FatalLevel

	Panic  = sublog.Panic
	Panicf = sublog.Panicf
	Debugf = sublog.Debugf
	Debug  = sublog.Debug

	Error  = sublog.Error
	Errorf = sublog.Errorf

	Info  = sublog.Info
	Infof = sublog.Infof
)

func ParseLevel(lvl string) (sublog.Level, error) {
	return sublog.ParseLevel(lvl)
}

func SetLevel(lv sublog.Level) {
	sublog.SetLevel(lv)
}

func TraceErrorf(format string, args ...interface{}) {
	debug.PrintStack()
	Errorf(format, args...)
}

func TraceError(args ...interface{}) {
	debug.PrintStack()
	Error(args...)
}
