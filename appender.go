package logx

// Appender accepts log entries
type Appender interface {

	// Append log line. Append should be thread-safe.
	Append(level, prefix, line string, tags ...string)
}
