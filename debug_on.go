// +build debug

package logx

import (
	"fmt"
)

var (
	debugBytes = []byte("DEBUG")
)

// Debug logs value with DEBUG severity level only
// if "debug" tag is provided on build.
func (l *Log) Debug(v ...interface{}) {
	l.output(debugBytes, []byte(fmt.Sprint(v...)))
}

// Debugf logs formatted value with DEBUG severity level only
// if "debug" tag is provided on build.
func (l *Log) Debugf(format string, v ...interface{}) {
	l.output(debugBytes, []byte(fmt.Sprintf(format, v...)))
}

// Debug logs value with DEBUG severity level only
// if "debug" tag is provided on build.
func Debug(v ...interface{}) {
	std.output(debugBytes, []byte(fmt.Sprint(v...)))
}

// Debugf logs formatted value with DEBUG severity level only
// if "debug" tag is provided on build.
func Debugf(format string, v ...interface{}) {
	std.output(debugBytes, []byte(fmt.Sprintf(format, v...)))
}
