package logger

import (
	"errors"
	"strings"

	"github.com/itfantasy/gonode/components"
	"github.com/itfantasy/gonode/components/rabbitmq"
	"github.com/itfantasy/gonode/utils/io"
	log "github.com/jeanphorn/log4go"
)

type Logger = log.Filter

var globalLogger log.Logger

func NewLogger(id string, loglevel string, logchan string, logcomp string) (*Logger, error) {
	var warn error = nil
	if strings.HasPrefix(logcomp, "rabbitmq://") {
		comp, err := components.NewComponent(logcomp)
		if err == nil {
			rmq, ok := comp.(*rabbitmq.RabbitMQ)
			if rmq != nil && ok {
				globalLogger = log.Logger{
					id: &log.Filter{getLogLevel(loglevel), NewRabbitMQLogWriter(rmq, logchan), id},
				}
				return globalLogger[id], nil
			} else {
				warn = errors.New("illegal log comp type! only rabbitmq or file or empty(console logger) ... ")
			}
		} else {
			warn = err
		}
	} else if strings.HasPrefix(logcomp, "file://") {
		filePath := strings.TrimPrefix(logcomp, "file://")
		if !io.FileExists(filePath) {
			dir := io.FetchDirByFilePath(filePath)
			io.MakeDir(dir)
		}
		globalLogger = log.Logger{
			id: &log.Filter{getLogLevel(loglevel), log.NewFileLogWriter(filePath, true, true), id},
		}
		return globalLogger[id], nil
	} else if logcomp != "" {
		warn = errors.New("illegal log comp type! only rabbitmq or file or empty(console logger) ... ")
	}
	globalLogger = log.Logger{
		id: &log.Filter{getLogLevel(loglevel), log.NewConsoleLogWriter(), id},
	}
	return globalLogger[id], warn
}

func getLogLevel(l string) log.Level {
	var lvl log.Level
	switch l {
	case "FINEST":
		lvl = log.FINEST
	case "FINE":
		lvl = log.FINE
	case "DEBUG":
		lvl = log.DEBUG
	case "TRACE":
		lvl = log.TRACE
	case "INFO":
		lvl = log.INFO
	case "WARNING":
		lvl = log.WARNING
	case "ERROR":
		lvl = log.ERROR
	case "CRITICAL":
		lvl = log.CRITICAL
	default:
		lvl = log.DEBUG
	}
	return lvl
}
