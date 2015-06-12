package logx

var std *Log = New()

func init() {
	std.calldepth = 3
}

func SetLevel(level string) {
	std.SetLevel(level)
}

func Trace(v ...interface{}) {
	std.Trace(v...)
}

func Tracef(format string, v ...interface{}) {
	std.Tracef(format, v...)
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

func ErrWarning(err interface{}, v ...interface{}) {
	std.ErrWarning(err, v...)
}

func ErrWarningf(err interface{}, format string, v ...interface{}) {
	std.ErrWarningf(err, format, v...)
}

func ErrError(err interface{}, v ...interface{}) {
	std.ErrError(err, v...)
}

func ErrErrorf(err interface{}, format string, v ...interface{}) {
	std.ErrErrorf(err, format, v...)
}

func ErrFatal(err interface{}, v ...interface{}) {
	std.ErrFatal(err, v...)
}

func ErrFatalf(err interface{}, format string, v ...interface{}) {
	std.ErrFatalf(err, format, v...)
}

func ErrPanic(err interface{}, v ...interface{}) {
	std.ErrPanic(err, v...)
}

func ErrPanicf(err interface{}, format string, v ...interface{}) {
	std.ErrPanicf(err, format, v...)
}

func prepend(head interface{}, tail []interface{}) []interface{} {
	return append([]interface{}{head}, tail...)
}