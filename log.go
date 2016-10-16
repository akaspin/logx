package logx

import (
	"sync/atomic"
	"unsafe"
	"fmt"
)

const (
	lInfo = "INFO"
	lWarning = "WARNING"
	lError = "ERROR"
	lCritical = "CRITICAL"

	// Ldate adds the date in the local time zone: 2009/01/23
	Ldate = 1 << iota

	// Ltime adds the time in the local time zone: 01:23:23
	Ltime

	// Lmicroseconds adds microsecond resolution: 01:23:23.123123.
	// assumes LTime.
	Lmicroseconds

	// Llongfile adds full file name and line number: /a/b/c/d.go:23
	Llongfile

	// Lshortfile adds final file name element and line number: d.go:23.
	// overrides Llongfile
	Lshortfile

	// LUTC if Ldate or Ltime is set, use UTC rather than the local time zone
	LUTC

	// LstdFlags initial values for the standard logger
	LstdFlags = Lshortfile
)

type Log struct {
	prefix   string
	tags []string

	appenderPtr *unsafe.Pointer
	callDepth int
}

// NewLog returns new log. Output Writer must be thread safe.
func NewLog(appender Appender, prefix string, tags ...string) (res *Log) {
	res = &Log{
		tags: tags,
		prefix: prefix,
		appenderPtr: new(unsafe.Pointer),
		callDepth: 2,
	}
	res.SetAppender(appender)
	return
}

// GetLog returns new independent log instance with given prefix.
func (l *Log) GetLog(prefix string, tags ...string) (res *Log) {
	res = NewLog(
		(*(*Appender)(atomic.LoadPointer(l.appenderPtr))),
		prefix, tags...)
	return
}

// SetAppender sets appender for Log instance.
func (l *Log) SetAppender(appender Appender) {
	atomic.StorePointer(l.appenderPtr, (unsafe.Pointer)(&appender))
}

// Prefix returns log prefix.
func (l *Log) Prefix() (res string) {
	res = string(l.prefix)
	return
}

// Tags returns log tags.
func (l *Log) Tags() (res []string) {
	res = l.tags
	return
}

// Print is synonym to Info used for compatibility.
func (l *Log) Print(v ...interface{}) {
	l.append(lInfo, fmt.Sprint(v...))
}

// Printf is synonym to Infof used for compatibility.
func (l *Log) Printf(format string, v ...interface{}) {
	l.append(lInfo, fmt.Sprintf(format, v...))
}

// Info logs value with INFO severity level.
func (l *Log) Info(v ...interface{}) {
	l.append(lInfo, fmt.Sprint(v...))
}

// Infof logs formatted value with INFO severity level.
func (l *Log) Infof(format string, v ...interface{}) {
	l.append(lInfo, fmt.Sprintf(format, v...))
}

// Warning logs value with WARNING severity level.
func (l *Log) Warning(v ...interface{}) {
	l.append(lWarning, fmt.Sprint(v...))
}

// Warningf logs formatted value with WARNING severity level.
func (l *Log) Warningf(format string, v ...interface{}) {
	l.append(lWarning, fmt.Sprintf(format, v...))
}

// Error logs value with ERROR severity level.
func (l *Log) Error(v ...interface{}) {
	l.append(lError, fmt.Sprint(v...))
}

// Errorf logs formatted value with ERROR severity level.
func (l *Log) Errorf(format string, v ...interface{}) {
	l.append(lError, fmt.Sprintf(format, v...))
}

// Critical logs value with CRITICAL severity level.
func (l *Log) Critical(v ...interface{}) {
	l.append(lCritical, fmt.Sprint(v...))
}

// Criticalf logs formatted value with CRITICAL severity level.
func (l *Log) Criticalf(format string, v ...interface{}) {
	l.append(lCritical, fmt.Sprintf(format, v...))
}

func (l *Log) append(level, line string) {
	(*(*Appender)(atomic.LoadPointer(l.appenderPtr))).Append(
		level, l.prefix, line, l.tags...)
}
