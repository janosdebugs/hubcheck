package hublog

import "fmt"

type Logger interface {
    Logf(format string, args ...interface{})
    Loge(err error)

    WithLevel(level Level) Logger
}

func New() Logger {
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
        prefix = "\033[37m⚙️ "
    case Info:
        prefix = "\033[36mℹ️️ "
    case Notice:
        prefix = "\033[32m✅️ "
    case Warning:
        prefix = "\033[33m⚠️️️ "
    case Error:
        prefix = "\033[31m❌️️ "
    }
    print(prefix + fmt.Sprintf(format, args...) + suffix)
}

func (l logger) Loge(err error) {
    l.Logf("%v", err)
}

func (l logger) WithLevel(level Level) Logger {
    return &logger{
        level: level,
    }
}
