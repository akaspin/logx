package logx

import (
	"sync/atomic"
	"io"
	"unsafe"
	"fmt"
	"time"
	"runtime"
)

const (
	INFO     = "INFO"
	WARNING  = "WARNING"
	ERROR    = "ERROR"
	CRITICAL = "CRITICAL"

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

var (
	infoBytes = []byte(INFO)
	warningBytes = []byte(WARNING)
	errorBytes = []byte(ERROR)
	criticalBytes = []byte(CRITICAL)
)

type Log struct {
	outputPtr *unsafe.Pointer
	flagsPtr  *int32
	prefix   []byte

	callDepth int
}

// NewLog returns new log. Output Writer must be thread safe.
func NewLog(output io.Writer, prefix string, flags int) (res *Log) {
	res = &Log{
		outputPtr: new(unsafe.Pointer),
		flagsPtr: new(int32),
		callDepth: 2,
	}
	if prefix != "" {
		res.prefix = []byte(prefix)
	}
	res.SetFlags(flags)
	res.SetOutput(output)
	return
}

// GetLog returns new independent log instance based on Log parameters.
func (l *Log) GetLog(prefix string) (res *Log) {
	res = NewLog(
		(*(*io.Writer)(atomic.LoadPointer(l.outputPtr))),
		prefix,
		l.Flags())
	return
}

// SetOutput thread-safely sets log output.
func (l *Log) SetOutput(output io.Writer) {
	atomic.StorePointer(l.outputPtr, (unsafe.Pointer)(&output))
}

// SetFlags thread-safely sets log flags.
func (l *Log) SetFlags(flags int) {
	atomic.StoreInt32(l.flagsPtr, int32(flags))
}

// Flags returns log flags.
func (l *Log) Flags() (res int) {
	res = int(atomic.LoadInt32(l.flagsPtr))
	return
}

// Prefix returns log prefix.
func (l *Log) Prefix() (res string) {
	res = string(l.prefix)
	return
}

// Print is synonym to Info used for compatibility.
func (l *Log) Print(v ...interface{}) {
	l.output(infoBytes, []byte(fmt.Sprint(v...)))
}

// Printf is synonym to Infof used for compatibility.
func (l *Log) Printf(format string, v ...interface{}) {
	l.output(infoBytes, []byte(fmt.Sprintf(format, v...)))
}

// Info logs value with INFO severity level.
func (l *Log) Info(v ...interface{}) {
	l.output(infoBytes, []byte(fmt.Sprint(v...)))
}

// Infof logs formatted value with INFO severity level.
func (l *Log) Infof(format string, v ...interface{}) {
	l.output(infoBytes, []byte(fmt.Sprintf(format, v...)))
}

// Warning logs value with WARNING severity level.
func (l *Log) Warning(v ...interface{}) {
	l.output(warningBytes, []byte(fmt.Sprint(v...)))
}

// Warningf logs formatted value with WARNING severity level.
func (l *Log) Warningf(format string, v ...interface{}) {
	l.output(warningBytes, []byte(fmt.Sprintf(format, v...)))
}

// Error logs value with ERROR severity level.
func (l *Log) Error(v ...interface{}) {
	l.output(errorBytes, []byte(fmt.Sprint(v...)))
}

// Errorf logs formatted value with ERROR severity level.
func (l *Log) Errorf(format string, v ...interface{}) {
	l.output(errorBytes, []byte(fmt.Sprintf(format, v...)))
}

// Critical logs value with CRITICAL severity level.
func (l *Log) Critical(v ...interface{}) {
		l.output(criticalBytes, []byte(fmt.Sprint(v...)))
}

// Criticalf logs formatted value with CRITICAL severity level.
func (l *Log) Criticalf(format string, v ...interface{}) {
		l.output(criticalBytes, []byte(fmt.Sprintf(format, v...)))
}

func (l *Log) output(level, line []byte) {
	var buf []byte
	flags := int(atomic.LoadInt32(l.flagsPtr))

	// time
	if flags&(Ldate|Ltime|Lmicroseconds|LUTC) != 0 {
		t := time.Now()
		if flags&LUTC != 0 {
			t = t.UTC()
		}
		if flags&(Ldate|Ltime|Lmicroseconds) != 0 {
			if flags&Ldate != 0 {
				year, month, day := t.Date()
				itoa(&buf, year, 4)
				buf = append(buf, '/')
				itoa(&buf, int(month), 2)
				buf = append(buf, '/')
				itoa(&buf, day, 2)
				buf = append(buf, ' ')
			}
			if flags&(Ltime|Lmicroseconds) != 0 {
				hour, min, sec := t.Clock()
				itoa(&buf, hour, 2)
				buf = append(buf, ':')
				itoa(&buf, min, 2)
				buf = append(buf, ':')
				itoa(&buf, sec, 2)
				if flags&Lmicroseconds != 0 {
					buf = append(buf, '.')
					itoa(&buf, t.Nanosecond()/1e3, 6)
				}
				buf = append(buf, ' ')
			}
		}
	}

	// level
	buf = append(buf, level...)
	buf = append(buf, ' ')

	// prefix
	if l.prefix != nil {
		buf = append(buf, l.prefix...)
		buf = append(buf, ' ')
	}

	// file
	if flags&(Lshortfile|Llongfile) != 0 {
		_, file, lineNo, ok := runtime.Caller(l.callDepth)
		if !ok {
			file = "???"
			lineNo = 0
		}
		if flags&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		buf = append(buf, file...)
		buf = append(buf, ':')
		itoa(&buf, lineNo, -1)
		buf = append(buf, ' ')
	}

	buf = append(buf, line...)

	if len(line) == 0 || line[len(line)-1] != '\n' {
		buf = append(buf, '\n')
	}

	(*(*io.Writer)(atomic.LoadPointer(l.outputPtr))).Write(buf)
}

// Cheap integer to fixed-width decimal ASCII.
// Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
