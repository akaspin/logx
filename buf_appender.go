package logx

import (
	"bytes"
	"io"
	"runtime"
	"sync"
	"time"
	"unicode"
)

// Default Appender based on buffer pool
type BufferAppender struct {
	output     io.Writer
	flags      int
	bufferPool sync.Pool
}

func NewBufferAppender(output io.Writer, flags int) (a *BufferAppender) {
	a = &BufferAppender{
		output: output,
		flags:  flags,
		bufferPool: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
	return
}

func (a *BufferAppender) Append(level, prefix, line string, tags ...string) {
	buf := a.bufferPool.Get().(*bytes.Buffer)

	// time
	if a.flags&(Ldate|Ltime|Lmicroseconds|LUTC) != 0 {
		t := time.Now()
		if a.flags&LUTC != 0 {
			t = t.UTC()
		}
		if a.flags&(Ldate|Ltime|Lmicroseconds) != 0 {
			if a.flags&Ldate != 0 {
				year, month, day := t.Date()
				itoaBuf(buf, year, 4)
				buf.WriteByte('/')
				itoaBuf(buf, int(month), 2)
				buf.WriteByte('/')
				itoaBuf(buf, day, 2)
				buf.WriteByte(' ')
			}
			if a.flags&(Ltime|Lmicroseconds) != 0 {
				hour, min, sec := t.Clock()
				itoaBuf(buf, hour, 2)
				buf.WriteByte(':')
				itoaBuf(buf, min, 2)
				buf.WriteByte(':')
				itoaBuf(buf, sec, 2)
				if a.flags&Lmicroseconds != 0 {
					buf.WriteByte('.')
					itoaBuf(buf, t.Nanosecond()/1e3, 6)
				}
				buf.WriteByte(' ')
			}
		}
	}

	// level
	buf.WriteString(level)
	buf.WriteByte(' ')

	// prefix
	if prefix != "" {
		buf.WriteString(prefix)
		buf.WriteByte(' ')
	}

	// tags
	ltags := len(tags)
	if ltags != 0 {
		buf.WriteByte('[')
		for i, tag := range tags {
			buf.WriteString(tag)
			if i < ltags-1 {
				buf.WriteByte(' ')
			}
		}
		buf.WriteString("] ")
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
		buf.WriteString(file)
		buf.WriteByte(':')
		itoaBuf(buf, lineNo, -1)
		buf.WriteByte(' ')
	}

	if a.flags&(Lcompact) != 0 {
		stripBuf(buf, line)
	} else {
		buf.WriteString(line)
	}

	if len(line) == 0 || line[len(line)-1] != '\n' {
		buf.WriteByte('\n')
	}
	buf.WriteTo(a.output)
	buf.Reset()
	a.bufferPool.Put(buf)
}

func itoaBuf(buf *bytes.Buffer, i int, wid int) {
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
	buf.Write(b[:])
}

func stripBuf(buf *bytes.Buffer, src string) {
	var white bool
	for _, c := range src {
		if unicode.IsSpace(c) {
			if !white {
				buf.WriteByte(' ')
			}
			white = true
		} else {
			buf.WriteRune(c)
			white = false
		}
	}
}
