package logx

import (
	"os"
	"io"
	"fmt"
)

var std = NewLog(os.Stderr, "", LstdFlags)

// GetLog returns new independent log instance based
// on standard Log parameters.
func GetLog(prefix string) *Log {
	return std.GetLog(prefix)
}

// SetOutput thread-safely sets standard Log output.
func SetOutput(output io.Writer) {
	std.SetOutput(output)
}

// SetFlags thread-safely sets standard log flags.
func SetFlags(flags int) {
	std.SetFlags(flags)
}

// Flags returns standard log flags.
func Flags() int {
	return std.Flags()
}

// Print is synonym to standard Log Info used for compatibility.
func Print(v ...interface{}) {
	std.output(infoBytes, []byte(fmt.Sprint(v...)))
}

// Printf is synonym to standard Log Infof used for compatibility.
func Printf(format string, v ...interface{}) {
	std.output(infoBytes, []byte(fmt.Sprintf(format, v...)))
}

// Info logs value with INFO severity level to standard Log.
func Info(v ...interface{}) {
	std.output(infoBytes, []byte(fmt.Sprint(v...)))
}

// Infof logs formatted value with INFO severity level to standard Log.
func Infof(format string, v ...interface{}) {
	std.output(infoBytes, []byte(fmt.Sprintf(format, v...)))
}

// Warning logs value with WARNING severity level to standard Log.
func Warning(v ...interface{}) {
	std.output(warningBytes, []byte(fmt.Sprint(v...)))
}

// Warningf logs formatted value with WARNING severity level to standard Log.
func Warningf(format string, v ...interface{}) {
	std.output(warningBytes, []byte(fmt.Sprintf(format, v...)))
}

// Error logs value with ERROR severity level to standard Log.
func Error(v ...interface{}) {
	std.output(errorBytes, []byte(fmt.Sprint(v...)))
}

// Errorf logs formatted value with ERROR severity level to standard Log.
func Errorf(format string, v ...interface{}) {
	std.output(errorBytes, []byte(fmt.Sprintf(format, v...)))
}

// Critical logs value with CRITICAL severity level to standard Log.
func Critical(v ...interface{}) {
	std.output(criticalBytes, []byte(fmt.Sprint(v...)))
}

// Criticalf logs formatted value with CRITICAL severity level to standard Log.
func Criticalf(format string, v ...interface{}) {
	std.output(criticalBytes, []byte(fmt.Sprintf(format, v...)))
}
