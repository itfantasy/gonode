package logger

import (
	"fmt"
	"time"

	"github.com/itfantasy/gonode/utils/json"
	"github.com/itfantasy/gonode/utils/ts"
)

type LogInfo struct {
	Level    int
	Created  string
	Source   string
	Message  string
	Category string
}

func ParseJson(str string) (*LogInfo, error) {
	l := new(LogInfo)
	err := json.Unmarshal(str, l)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (l *LogInfo) SetCreated(tm time.Time) {
	l.Created = tm.Format(ts.FORMAT_NOW_A)
}

func (l *LogInfo) ToJson() (string, error) {
	return json.Marshal(l)
}

func (l *LogInfo) FormatString() string {
	return fmt.Sprintf("[%s] [%s] [%s] (%s) %s",
		l.Created,
		l.Category,
		LevelToString(l.Level),
		l.Source,
		l.Message)
}

func (l *LogInfo) Println() {
	fmt.Println(l.FormatString())
}
