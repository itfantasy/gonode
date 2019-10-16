package logger

import (
	"strconv"
	"strings"
)

const (
	DEBUG int = 0
	INFO      = 100
	WARN      = 500
	ERROR     = 1000
	FATAL     = 9999
)

func LevelToString(lv int) string {
	if lv < INFO {
		return "DEBUG"
	} else if lv == INFO {
		return "INFO"
	} else if lv > INFO && lv < WARN {
		return "INFO-lv" + strconv.Itoa(lv-INFO)
	} else if lv == WARN {
		return "WARN"
	} else if lv > WARN && lv < ERROR {
		return "WARN-lv" + strconv.Itoa(lv-WARN)
	} else if lv == ERROR {
		return "ERROR"
	} else if lv > ERROR && lv < FATAL {
		return "ERROR-lv" + strconv.Itoa(lv-ERROR)
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
