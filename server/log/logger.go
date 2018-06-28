package log

import "context"

// Level is the message urgency level.
type Level uint

const (
	_ = iota
	// LevelDebug used for debug messages.
	LevelDebug
	// LevelInfo used for info messages.
	LevelInfo
	// LevelWarning used for warning messages.
	LevelWarning
	// LevelError used for error messages.
	LevelError
	// LevelFatal used for fatal messages. os.Exit(1) is called after printing.
	LevelFatal
)

// Writer accepts messages along with the Level.
type Writer interface {
	Log(Level, ...interface{})
	Logf(Level, string, ...interface{})
}

// WriterC is the context-aware Writer.
type WriterC interface {
	Logc(context.Context, Level, ...interface{})
	Logcf(context.Context, Level, string, ...interface{})
}

// Default is the default logger.
var Default Writer
