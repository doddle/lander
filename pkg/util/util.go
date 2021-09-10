package util

type LoggerIFace interface {
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
}
