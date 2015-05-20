/*
Simple wrapper around log.
*/

package logx

import (
	"log"
	"os"
	"strings"
	"runtime"
	"path/filepath"
	"fmt"
)

const (
	debugLog = 0
	infoLog = 1
	warningLog = 2
	errorLog = 3
	fatalLog = 4
	panicLog = 5

	DEBUG = "DEBUG"
	INFO = "INFO"
	WARNING = "WARNING"
	ERROR = "ERROR"
	FATAL = "FATAL"
	PANIC = "PANIC"
)

var (
	logLevels = []string{
		"DEBUG",
		"INFO",
		"WARNING",
		"ERROR",
		"FATAL",
		"PANIC",
	}
	std *Log = New()
)

func init() {
	std.calldepth = 3
}

// Logger
type Log struct  {
	logger *log.Logger

	level int
	calldepth int
}

// Get new logger
func New() *Log {
	return &Log{
		logger: log.New(os.Stderr, "", 0),
		calldepth: 2,
		level: errorLog,
	}
}

// Set logging level: DEBUG, INFO, WARNING, ERROR, FATAL, PANIC
func (l *Log) SetLevel(level string) {
	level = strings.ToUpper(level)
	for i, st := range logLevels {
		if st == level {
			l.level = i
			return
		}
	}
	l.level = errorLog
}

func (l *Log) Debug(v ...interface{}) {
	l.output(debugLog, l.header(DEBUG, fmt.Sprint(v...)))
}

func (l *Log) Debugf(format string, v ...interface{}) {
	l.output(debugLog, l.header(DEBUG, fmt.Sprintf(format, v...)))
}

func (l *Log) Info(v ...interface{}) {
	l.output(infoLog, l.header(INFO, fmt.Sprint(v...)))
}

func (l *Log) Infof(format string, v ...interface{}) {
	l.output(infoLog, l.header(INFO, fmt.Sprintf(format, v...)))
}

func (l *Log) Warning(v ...interface{}) {
	l.output(warningLog, l.header(WARNING, fmt.Sprint(v...)))
}

func (l *Log) Warningf(format string, v ...interface{}) {
	l.output(warningLog, l.header(WARNING, fmt.Sprintf(format, v...)))
}

func (l *Log) Error(v ...interface{}) {
	l.output(errorLog, l.header(ERROR, fmt.Sprint(v...)))
}

func (l *Log) Errorf(format string, v ...interface{}) {
	l.output(errorLog, l.header(ERROR, fmt.Sprintf(format, v...)))
}

func (l *Log) Fatal(v ...interface{}) {
	l.output(fatalLog, l.header(FATAL, fmt.Sprint(v...)))
	os.Exit(1)
}

func (l *Log) Fatalf(format string, v ...interface{}) {
	l.output(fatalLog, l.header(ERROR, fmt.Sprintf(format, v...)))
	os.Exit(1)
}

func (l *Log) Panic(v ...interface{}) {
	msg := l.header(PANIC, fmt.Sprint(v...))
	l.output(panicLog, msg)
	panic(msg)
}

func (l *Log) Panicf(format string, v ...interface{}) {
	msg := l.header(PANIC, fmt.Sprintf(format, v...))
	l.output(panicLog, msg)
	panic(msg)
}

func (l *Log) output(lvl int, line string) {
	if lvl >= l.level {
		l.logger.Output(l.calldepth, line)
	}
}

func (l *Log) header(lvl, msg string) string {
	_, file, line, ok := runtime.Caller(l.calldepth)
	if ok {
		file = filepath.Base(file)
	}

	if len(file) == 0 {
		file = "???"
	}

	if line < 0 {
		line = 0
	}

	return fmt.Sprintf("%s %s:%d: %s", lvl, file, line, msg)
}

func SetLevel(level string) {
	std.SetLevel(level)
}

func Debug(v ...interface{}) {
	std.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	std.Debugf(format, v...)
}

func Info(v ...interface{}) {
	std.Info(v...)
}

func Infof(format string, v ...interface{}) {
	std.Infof(format, v...)
}

func Warning(v ...interface{}) {
	std.Warning(v...)
}

func Warningf(format string, v ...interface{}) {
	std.Warningf(format, v...)
}

func Error(v ...interface{}) {
	std.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	std.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	std.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	std.Fatalf(format, v...)
}

func Panic(v ...interface{}) {
	std.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	std.Panicf(format, v...)
}
