package log

import (
	"fmt"

	"github.com/sirupsen/logrus"
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
		logrus.Debug(i...)
	case LevelInfo:
		logrus.Info(i...)
	case LevelWarning:
		logrus.Warn(i...)
	case LevelError:
		logrus.Error(i...)
	case LevelFatal:
		logrus.Fatal(i...)
	}
}

// Logf implements Dumper.
func (s Logrus) Logf(l Level, msg string, args ...interface{}) {
	s.Log(l, fmt.Sprintf(msg, args...))
}
