package hublog

import (
	"fmt"
	"os"
)

type Logger interface {
	Logf(format string, args ...interface{})
	Loge(err error)

	WithLevel(level Level) Logger
}

func New(minLogLevel Level) Logger {

	switch minLogLevel {
	case Debug:
	case Info:
	case Notice:
	case Warning:
	case Error:
	default:
		panic("Invalid log level: " + minLogLevel)
	}

	return &logger{}
}

type logger struct {
	level Level
}

func (l logger) Logf(format string, args ...interface{}) {
	prefix := ""
	suffix := "\033[0m\n"
	switch l.level {
	case Debug:
		if l.level != Debug {
			return
		}
		prefix = "\033[37m⚙️ "
	case Info:
		if l.level != Debug && l.level != Info {
			return
		}
		prefix = "\033[36mℹ️️️ "
	case Notice:
		if l.level != Debug && l.level != Info && l.level != Notice {
			return
		}
		prefix = "\033[32m✅️ "
	case Warning:
		if l.level != Debug && l.level != Info && l.level != Notice && l.level != Warning {
			return
		}
		prefix = "\033[33m⚠️️️ "
	case Error:
		prefix = "\033[31m❌️️ "
	}
	_, _ = os.Stderr.WriteString(prefix + fmt.Sprintf(format, args...) + suffix)
}

func (l logger) Loge(err error) {
	l.Logf("%v", err)
}

func (l logger) WithLevel(level Level) Logger {
	return &logger{
		level: level,
	}
}
