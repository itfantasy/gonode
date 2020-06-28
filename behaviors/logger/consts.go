package logger

import (
	"strings"
)

const (
	DEBUG int = 0
	INFO      = 1
	WARN      = 2
	ERROR     = 3
	FATAL     = 4
)

func LevelToString(lv int) string {
	if lv == DEBUG {
		return "DEBUG"
	} else if lv == INFO {
		return "INFO"
	} else if lv == WARN {
		return "WARN"
	} else if lv == ERROR {
		return "ERROR"
	} else {
		return "FATAL"
	}
}

func StringToLevel(lv string) int {
	infos := strings.Split(lv, "-")
	if len(infos) <= 0 {
		return DEBUG
	}
	var lvl int = DEBUG
	switch infos[0] {
	case "DEBUG":
		lvl = DEBUG
	case "INFO":
		lvl = INFO
	case "WARNING":
		lvl = WARN
	case "ERROR":
		lvl = ERROR
	case "FATAL":
		lvl = FATAL
	default:
		lvl = DEBUG
	}
	return lvl
}
