package cache

import (
	"os"
	log "log"
	"fmt"
)

const (
	// Default name of default logger
	Default = "default"
)

var enabled = false
var lmap = map[string]Logger{}
var std = log.New(os.Stderr, "", log.Ldate | log.Ltime | log.Lshortfile)

// Logger interface
type Logger interface {
	Trace(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Critical(format string, v ...interface{})
}

// EnableLogger enable logger
// param enable - true is enable
func EnableLogger(enable bool) {
	enabled = enable
}

// SetLogger set logger
// param logger - instance of logger
// return arg1 - Error
func SetLogger(logger Logger) error {
	enabled = true
	lmap["default"] = logger
	return nil
}

// Trace output trace log
func Trace(format string, v ...interface{}) {
	if (!enabled) {
		return
	}
	if lmap[Default] != nil {
		lmap[Default].Trace(format, v...)
	} else {
		std.Output(2, fmt.Sprintf(format, v...))
	}
}

// Debug output debug log
func Debug(format string, v ...interface{}) {
	if (!enabled) {
		return
	}
	if lmap[Default] != nil {
		lmap[Default].Debug(format, v...)
	} else {
		std.Output(2, fmt.Sprintf(format, v...))
	}
}

// Info output information log
func Info(format string, v ...interface{}) {
	if (!enabled) {
		return
	}
	if lmap[Default] != nil {
		lmap[Default].Info(format, v...)
	} else {
		std.Output(2, fmt.Sprintf(format, v...))
	}
}

// Warn output warning log
func Warn(format string, v ...interface{}) {
	if (!enabled) {
		return
	}
	if lmap[Default] != nil {
		lmap[Default].Warn(format, v...)
	} else {
		std.Output(2, fmt.Sprintf(format, v...))
	}
}

// Error output error log
func Error(format string, v ...interface{}) {
	if (!enabled) {
		return
	}
	if lmap[Default] != nil {
		lmap[Default].Error(format, v...)
	} else {
		std.Output(2, fmt.Sprintf(format, v...))
	}
}

// Critical output critical log
func Critical(format string, v ...interface{}) {
	if (!enabled) {
		return
	}
	if lmap[Default] != nil {
		lmap[Default].Critical(format, v...)
	} else {
		std.Output(2, fmt.Sprintf(format, v...))
	}
}
