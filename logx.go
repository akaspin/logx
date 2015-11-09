/*
Simple wrapper around log.
*/

package logx

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"io"
)

const (
	traceLog   = 0
	debugLog   = 1
	infoLog    = 2
	warningLog = 3
	errorLog   = 4
	fatalLog   = 5
	panicLog   = 6

	TRACE   = "TRACE"
	DEBUG   = "DEBUG"
	INFO    = "INFO"
	WARNING = "WARNING"
	ERROR   = "ERROR"
	FATAL   = "FATAL"
	PANIC   = "PANIC"
)

var (
	logLevels = []string{
		"TRACE",
		"DEBUG",
		"INFO",
		"WARNING",
		"ERROR",
		"FATAL",
		"PANIC",
	}
)

type Logger interface {
	SetLevel(level string)
	SetOutput(w io.Writer)
	SetPrefix(prefix string)

	Trace(v ...interface{})
	Tracef(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warning(v ...interface{})
	Warningf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})

	OnWarning(err interface{}, v ...interface{}) error
	OnWarningf(err interface{}, format string, v ...interface{})
	OnError(err interface{}, v ...interface{}) error
	OnErrorf(err interface{}, format string, v ...interface{}) error
	OnFatal(err interface{}, v ...interface{})
	OnFatalf(err interface{}, format string, v ...interface{})
	OnPanic(err interface{}, v ...interface{})
}

// Logger
type Log struct {
	logger *log.Logger

	level     int
	calldepth int
}

// Get new logger
func New() *Log {
	return &Log{
		logger:    log.New(os.Stderr, "", 0),
		calldepth: 2,
		level:     errorLog,
	}
}

// Set logging level: TRACE, DEBUG, INFO, WARNING, ERROR, FATAL, PANIC
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

// Set output
func (l *Log) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

func (l *Log) SetPrefix(prefix string) {
	l.logger.SetPrefix(prefix)
}

func (l *Log) Trace(v ...interface{}) {
	l.output(traceLog, l.header(TRACE, fmt.Sprint(v...)))
}

func (l *Log) Tracef(format string, v ...interface{}) {
	l.output(traceLog, l.header(TRACE, fmt.Sprintf(format, v...)))
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

// Acts as Warning(err, v...) but only if err argument is not nil
// OnWarning always returns <nil>. Use this for reports.
func (l *Log) OnWarning(err interface{}, v ...interface{}) error {
	if err != nil {
		l.output(warningLog, l.header(WARNING, fmt.Sprint(prepend(err, v)...)))
	}
	return nil
}

// Acts as Warningf(format, v...) but only if err is not nil
func (l *Log) OnWarningf(err interface{}, format string, v ...interface{}) error {
	if err != nil {
		l.output(warningLog, l.header(WARNING, fmt.Sprintf(format, v...)))
	}
	return nil
}

// Acts as Error(err, v...) but only if err argument is not nil
func (l *Log) OnError(err interface{}, v ...interface{}) error {
	if err != nil {
		l.output(errorLog, l.header(ERROR, fmt.Sprint(prepend(err, v)...)))
	}
	return nil
}

// Acts as Errorf(format, v...) but only if err is not nil
func (l *Log) OnErrorf(err interface{}, format string, v ...interface{}) error {
	if err != nil {
		l.output(errorLog, l.header(ERROR, fmt.Sprintf(format, v...)))
	}
	return nil
}

// Acts as Fatal(err, v...) but only if err argument is not nil
func (l *Log) OnFatal(err interface{}, v ...interface{}) {
	if err != nil {
		l.output(fatalLog, l.header(FATAL, fmt.Sprint(prepend(err, v)...)))
		os.Exit(1)
	}
}

// Acts as Fatalf(format, err, v...) but only if err is not nil
func (l *Log) OnFatalf(err interface{}, format string, v ...interface{}) {
	if err != nil {
		l.output(fatalLog, l.header(ERROR, fmt.Sprintf(format, v...)))
		os.Exit(1)
	}
}

// Acts as Panic() but only if err argument is not nil
func (l *Log) OnPanic(err interface{}, v ...interface{}) {
	if err != nil {
		msg := l.header(PANIC, fmt.Sprint(prepend(err, v)...))
		l.output(panicLog, msg)
		panic(msg)
	}
}

// Acts as Panicf(format, err, v...) but only if err is not nil
func (l *Log) OnPanicf(err interface{}, format string, v ...interface{}) {
	if err != nil {
		msg := l.header(PANIC, fmt.Sprintf(format, v...))
		l.output(panicLog, msg)
		panic(msg)
	}
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

func prepend(head interface{}, tail []interface{}) []interface{} {
	return append([]interface{}{head}, tail...)
}
