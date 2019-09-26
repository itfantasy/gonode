package logger

import (
	"github.com/itfantasy/gonode/utils/json"
)

type LogInfo struct {
	Level    int
	Created  string
	Source   string
	Message  string
	Category string
}

func ParseLog(str string) (*LogInfo, error) {
	l := new(LogInfo)
	err := json.Unmarshal(str, l)
	if err != nil {
		return nil, err
	}
	return l, nil
}
