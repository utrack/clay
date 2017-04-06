package log

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

// Dumper accepts messages along with the Level.
type Dumper interface {
	Log(Level, ...interface{})
	Logf(Level, string, ...interface{})
}

// Default is the default logger.
var Default Dumper
