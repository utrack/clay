package log

import (
	"fmt"
	"os"
	"time"
)

func init() {
	Default = Logrus{}
}

// Logrus is the default logger, using sirupsen/logrus as a backend.
type Logrus struct{}

// Log implements Dumper.
func (s Logrus) Log(l Level, i ...interface{}) {
	switch l {
	case LevelDebug:
		fallthrough
	case LevelInfo:
		fallthrough
	case LevelWarning:
		fallthrough
	case LevelError:
		fallthrough
	case LevelFatal:
		s.log(l, i...)
	}
}

// Logf implements Dumper.
func (s Logrus) Logf(l Level, msg string, args ...interface{}) {
	s.Log(l, fmt.Sprintf(msg, args...))
}

// Default logger implementation. For backward compatibility it is taken from
// @link github.com/sirupsen/logrus@v1.0.5/text_formatter.go
func (s Logrus) log(level Level, args ...interface{}) {
	fmt.Println("time=" + time.Now().Format(time.RFC3339) +
		" level=" + levelToString(level) +
		" msg=\"" + fmt.Sprint(args...) + "\"")

	if level == LevelFatal {
		os.Exit(1)
	}
}

// Convert the Level to a string.
// @link github.com/sirupsen/logrus@v1.0.5/logrus.go
func levelToString(l Level) string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarning:
		return "warning"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	}

	return "unknown"
}
