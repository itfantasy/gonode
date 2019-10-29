package logger

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

type Logger struct {
	category  string
	level     int
	logWriter LogWriter
}

func NewLogger(category string, loglevel string, logcomp string, logchan string) (*Logger, error) {
	l := new(Logger)
	l.category = category
	l.level = StringToLevel(loglevel)
	var warn error = nil
	if strings.HasPrefix(logcomp, "rabbitmq://") {
		if logWriter, err := NewRabbitMQLogWriter(logcomp, logchan); err != nil {
			warn = err
		} else {
			l.logWriter = logWriter
		}
	} else if strings.HasPrefix(logcomp, "file://") {
		filePath := strings.TrimPrefix(logcomp, "file://")
		if logWriter, err := NewFileLogWriter(filePath); err != nil {
			warn = err
		} else {
			l.logWriter = logWriter
		}
	} else {
		l.logWriter = NewConsoleLogWriter()
	}
	return l, warn
}

func (log *Logger) Source(callstack int) string {
	pc, _, lineno, ok := runtime.Caller(callstack)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}
	return src
}

func (log *Logger) Log4Extend(lvl int, callstack int, any interface{}, args ...interface{}) {
	if lvl < log.level {
		return
	}
	src := log.Source(callstack)
	var msg string = ""
	switch any.(type) {
	case string:
		msg = any.(string)
		if len(args) > 0 {
			msg = fmt.Sprintf(msg, args...)
		}
	case error:
		msg = any.(error).Error()
		if len(args) > 0 {
			msg = fmt.Sprintf(msg, args...)
		}
	default:
		msg = fmt.Sprint(any)
	}
	info := new(LogInfo)
	info.Category = log.category
	info.Level = lvl
	info.Message = msg
	info.Source = src
	info.SetCreated(time.Now())
	if lvl <= DEBUG {
		info.Println() // DEBUG always only to console
	} else {
		log.logWriter.LogWrite(info)
	}
}

func (log *Logger) Log(lvl int, arg0 interface{}, args ...interface{}) {
	log.Log4Extend(lvl, 2, arg0, args...)
}

func (log *Logger) Debug(arg0 interface{}, args ...interface{}) {
	log.Log4Extend(DEBUG, 2, arg0, args...)
}

func (log *Logger) Info(arg0 interface{}, args ...interface{}) {
	log.Log4Extend(INFO, 2, arg0, args...)
}

func (log *Logger) Warn(arg0 interface{}, args ...interface{}) {
	log.Log4Extend(WARN, 2, arg0, args...)
}

func (log *Logger) Error(arg0 interface{}, args ...interface{}) {
	log.Log4Extend(ERROR, 2, arg0, args...)
}

func (log *Logger) Fatal(arg0 interface{}, args ...interface{}) {
	log.Log4Extend(FATAL, 2, arg0, args...)
}
