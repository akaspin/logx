package logx

import (
	"io"
	"time"
	"runtime"
)

type Appender interface {

	// Append log line
	Append(level, prefix, line string, tags ...string)
}

type SimpleAppender struct {
	output io.Writer
	flags int
}

func NewSimpleAppender(output io.Writer, flags int) *SimpleAppender {
	return &SimpleAppender{
		output: output,
		flags: flags,
	}
}

func (a *SimpleAppender) Append(level, prefix, line string, tags ...string) {
	var buf []byte

	// time
	if a.flags&(Ldate|Ltime|Lmicroseconds|LUTC) != 0 {
		t := time.Now()
		if a.flags&LUTC != 0 {
			t = t.UTC()
		}
		if a.flags&(Ldate|Ltime|Lmicroseconds) != 0 {
			if a.flags&Ldate != 0 {
				year, month, day := t.Date()
				itoa(&buf, year, 4)
				buf = append(buf, '/')
				itoa(&buf, int(month), 2)
				buf = append(buf, '/')
				itoa(&buf, day, 2)
				buf = append(buf, ' ')
			}
			if a.flags&(Ltime|Lmicroseconds) != 0 {
				hour, min, sec := t.Clock()
				itoa(&buf, hour, 2)
				buf = append(buf, ':')
				itoa(&buf, min, 2)
				buf = append(buf, ':')
				itoa(&buf, sec, 2)
				if a.flags&Lmicroseconds != 0 {
					buf = append(buf, '.')
					itoa(&buf, t.Nanosecond()/1e3, 6)
				}
				buf = append(buf, ' ')
			}
		}
	}

	// level
	buf = append(buf, []byte(level)...)
	buf = append(buf, ' ')

	// prefix
	if prefix != "" {
		buf = append(buf, []byte(prefix)...)
		buf = append(buf, ' ')
	}

	// tags
	ltags := len(tags)
	if ltags > 0 {
		buf = append(buf, '[')
		for i, tag := range tags {
			buf = append(buf, tag...)
			if i < ltags-1 {
				buf = append(buf, ' ')
			}
		}
		buf = append(buf, "] "...)
	}

	// file
	if a.flags&(Lshortfile|Llongfile) != 0 {
		_, file, lineNo, ok := runtime.Caller(3)
		if !ok {
			file = "???"
			lineNo = 0
		}
		if a.flags&Lshortfile != 0 {
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

	buf = append(buf, []byte(line)...)

	if len(line) == 0 || line[len(line)-1] != '\n' {
		buf = append(buf, '\n')
	}
	a.output.Write(buf)
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
