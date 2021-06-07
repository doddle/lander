package main

import (
	"github.com/withmandala/go-log"
	"os"
)

func newLogger(debug bool) *log.Logger {
	// check if debug enabled
	if debug {
		logger := log.New(os.Stdout).WithDebug().WithColor()
		return logger
	}
	return log.New(os.Stdout).WithColor()
}
