package logger

import (
	log "github.com/jeanphorn/log4go"
)

var globalLogger log.Logger

func NewLogger(id string, loglevel string, rmqurl string, rmqhost string, rmqchan string, rmquser string, rmqpass string) *log.Filter {
	if rmqurl != "" && rmqhost != "" {
		globalLogger = log.Logger{
			id: &log.Filter{getLogLevel(loglevel), NewRabbitMQLogWriter(rmqurl, rmqhost, rmqchan, rmquser, rmqpass), id},
		}
	} else {
		globalLogger = log.Logger{
			id: &log.Filter{getLogLevel(loglevel), log.NewConsoleLogWriter(), id},
		}
	}
	return globalLogger[id]
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
