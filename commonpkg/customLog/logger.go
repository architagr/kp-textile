package customLog

import (
	"fmt"
	"io"
	"log"
)

type ILogger interface {
	Write(level int, message string)
}

type Logger struct {
	level   int
	service string
	l       *log.Logger
}

func (l *Logger) Write(level int, message string) {
	if level >= l.level {
		l.l.Print(message)
	}
}

var logger ILogger

func Init(le int, svcName string, out io.Writer) (ILogger, error) {
	var err error
	if logger == nil {
		if le <= 5 && le > 0 {
			logger = &Logger{
				level:   le,
				service: svcName,
				l:       log.New(out, svcName, log.Lmsgprefix|log.Ldate|log.LUTC|log.Lshortfile),
			}
		} else {
			return nil, fmt.Errorf("Level should be between 1 to 5")
		}
	}
	return logger, err
}
