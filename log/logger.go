package log

import (
	"github.com/utrack/clay/v2/server/log"
)

// Logrus is the default logger, using sirupsen/logrus as a backend.
type Logrus = log.Logrus

// Level is the message urgency level.
type Level log.Level

const (
	// LevelDebug used for debug messages.
	LevelDebug = log.LevelDebug
	// LevelInfo used for info messages.
	LevelInfo = log.LevelInfo
	// LevelWarning used for warning messages.
	LevelWarning = log.LevelWarning
	// LevelError used for error messages.
	LevelError = log.LevelError
	// LevelFatal used for fatal messages. os.Exit(1) is called after printing.
	LevelFatal = log.LevelFatal
)

// Writer accepts messages along with the Level.
type Writer = log.Writer

// WriterC is the context-aware Writer.
type WriterC = log.WriterC

// Default is the default logger.
var Default Writer = log.Default
