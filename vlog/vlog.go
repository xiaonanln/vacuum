package vlog

import (
	"runtime/debug"

	sublog "github.com/Sirupsen/logrus"
)

var (
	Debugf = sublog.Debugf
)

func Errorf(format string, args ...interface{}) {
	sublog.Errorf(format, args...)
}

func TraceErrorf(format string, args ...interface{}) {
	debug.PrintStack()
	sublog.Errorf(format, args...)
}
