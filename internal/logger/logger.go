// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package logger

import "log"

const (
	debug = "[DEBUG] "
	info  = "[INFO] "
	warn  = "[WARN] "
	err   = "[ERROR] "
)

func Debugf(format string, v ...interface{}) {
	loggerPrintf(debug, format, v)
}

func Debug(v ...interface{}) {
	loggerPrintf(debug, "%v", v)
}

func Infof(format string, v ...interface{}) {
	loggerPrintf(info, format, v)
}

func Info(v ...interface{}) {
	loggerPrintf(info, "%v", v)
}

func Warnf(format string, v ...interface{}) {
	loggerPrintf(warn, format, v)
}

func Warn(v ...interface{}) {
	loggerPrintf(warn, "%v", v)
}

func Errorf(format string, v ...interface{}) {
	loggerPrintf(err, format, v)
}

func Error(v ...interface{}) {
	loggerPrintf(err, "%v", v)
}

func loggerPrintf(level, format string, v ...interface{}) {
	log.Printf(level+format, v...)
}
