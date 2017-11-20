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

func Now() int64 {
	return time.Now().Unix()
}

func Day() int64 {
	return NowToDay(Now())
}

func StrToNow(str string, format string) int64 {
	tm, err := time.Parse(format, str)
	if err != nil {
		return 0
	}
	return tm.Unix()
}

func NowToStr(now int64, format string) string {
	tm := time.Unix(now, 0)
	return tm.Format(format)
}

func StrToDay(str string, format string) int64 {
	if str == "" {
		return Day()
	}

	var now int64
	tm, err := time.Parse(format, str)
	if err == nil {
		now = tm.Unix()
	}

	return NowToDay(now)
}

func DayToStr(day int64, format string) string {
	now := DayToNow(day)
	tm := time.Unix(now, 0)
	return tm.Format(format)
}

func NowToDay(now int64) int64 {
	return (now + 28800) / 86400
}

func DayToNow(day int64) int64 {
	return day*86400 - 28800
}
