package logger

import (
	"fmt"
)

type Logger struct {
	level       Enum_LogLevels
	netReporter INetReporter
}

// the log levels
type Enum_LogLevels int

const (
	ALL Enum_LogLevels = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

func (this *Logger) SetLevel(level Enum_LogLevels) {
	this.level = level
}

func (this *Logger) Level() Enum_LogLevels {
	return this.level
}

func (this *Logger) SetNetReporter(netReporter INetReporter) {
	this.netReporter = netReporter
}

func (this *Logger) Log(level Enum_LogLevels, msg string) {
	if this.level > level {
		return
	}
	txt := this.levelToString(level) + msg
	fmt.Println(txt)
	// call io to write log file
	//....

	if this.netReporter != nil {
		// call the interface to report the net log
		this.netReporter.ReportLog(txt)
	}
}

func (this *Logger) Debug(msg string) {
	this.Log(DEBUG, msg)
}

func (this *Logger) Info(msg string) {
	this.Log(INFO, msg)
}

func (this *Logger) Warn(msg string) {
	this.Log(WARN, msg)
}

func (this *Logger) Error(msg string) {
	this.Log(ERROR, msg)
}

func (this *Logger) Fatal(msg string) {
	this.Log(FATAL, msg)
}

func (this *Logger) levelToString(level Enum_LogLevels) string {
	switch level {
	case DEBUG:
		return "d:/"
	case INFO:
		return "i:/"
	case WARN:
		return "w:/"
	case ERROR:
		return "e:/"
	case FATAL:
		return "f:/"
	}
	return ""
}
