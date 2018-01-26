package logx

// Appender accepts log entries
type Appender interface {
	Append(level, line string)
	Clone(prefix string, tags []string) Appender
}
