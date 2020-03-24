package ts

import (
	"time"
)

// yyyy-MM-dd hh:mm:ss
const FORMAT_NOW_A = "2006-01-02 15:04:05"

// MM/dd/yyyy hh:mm:ss PM
const FORMAT_NOW_B = "01/02/2006 15:04:05"

// yyyyMMddhhmmss
const FORMAT_NOW_C = "20060102150405"

// yyyy-MM-dd
const FORMAT_DAY_A = "2006-01-02"

// MM/dd/yyyy
const FORMAT_DAY_B = "01/02/2006"

// yyyyMMdd
const FORMAT_DAY_C = "20060102"

func Time() int64 {
	return time.Now().Unix()
}

func MilliSecond() int64 {
	return time.Now().UnixNano() / 1e6
}

func Day() int64 {
	return TimeToDay(Time())
}

func StrToTime(str string, format string) int64 {
	tm, err := time.Parse(format, str)
	if err != nil {
		return 0
	}
	return tm.Unix()
}

func TimeToStr(sec int64, format string) string {
	tm := time.Unix(sec, 0)
	return tm.Format(format)
}

func StrToDay(str string, format string) int64 {
	if str == "" {
		return Day()
	}

	var sec int64
	tm, err := time.Parse(format, str)
	if err == nil {
		sec = tm.Unix()
	}

	return TimeToDay(sec)
}

func DayToStr(day int64, format string) string {
	sec := DayToTime(day)
	tm := time.Unix(sec, 0)
	return tm.Format(format)
}

func TimeToDay(sec int64) int64 {
	return (sec + 28800) / 86400
}

func DayToTime(day int64) int64 {
	return day*86400 - 28800
}
