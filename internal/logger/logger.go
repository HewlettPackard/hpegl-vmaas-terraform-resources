package logger

import "log"

type Logger struct{}

const (
	debug = "[DEBUG] "
	info  = "[INFO] "
	warn  = "[WARN] "
	err   = "[ERROR] "
)

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.printf(debug, format, v)
}

func (l *Logger) Debug(v ...interface{}) {
	l.printf(debug, "%v", v)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.printf(info, format, v)
}

func (l *Logger) Info(v ...interface{}) {
	l.printf(info, "%v", v)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.printf(warn, format, v)
}

func (l *Logger) Warn(v ...interface{}) {
	l.printf(warn, "%v", v)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.printf(err, format, v)
}

func (l *Logger) Error(v ...interface{}) {
	l.printf(err, "%v", v)
}

func (l *Logger) printf(level, format string, v ...interface{}) {
	log.Printf(level+format, v...)
}
